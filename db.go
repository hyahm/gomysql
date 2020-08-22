package gomysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hyahm/golog"
)

type Db struct {
	conn    *sql.DB
	conf    string
	Ctx     context.Context
	sql     string
	debug   bool
	sc      *Sqlconfig
	f       *os.File
	mu      *sync.RWMutex
	maxConn int
}

func (d *Db) conndb() error {
	conn, err := sql.Open("mysql", d.conf)
	if err != nil {
		golog.Info(err)
		return err
	}
	if err = conn.Ping(); err != nil {
		golog.Info(err)
		return err
	}
	d.conn = conn
	if d.sc.ReadTimeout == 0 {
		d.sc.ReadTimeout = time.Second * 30
	}
	if d.sc.MaxIdleConns > 0 {
		d.conn.SetMaxIdleConns(d.sc.MaxIdleConns)
	}
	if d.sc.MaxOpenConns > 0 {
		d.conn.SetMaxOpenConns(d.sc.MaxOpenConns)
		d.maxConn = d.sc.MaxOpenConns
	} else {
		d.maxConn = 1024
	}
	// 防止开始就有很多连接，导致

	d.conn.SetConnMaxLifetime(d.sc.ConnMaxLifetime)
	if d.sc.WriteLogWhenFailed {
		d.mu = &sync.RWMutex{}
		if d.sc.LogFile == "" {
			d.sc.LogFile = ".failed.sql"
		}
		var err error
		d.f, err = os.OpenFile(d.sc.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (d *Db) execError(err error, cmd string, args ...interface{}) (int64, error) {
	if d.sc.WriteLogWhenFailed {
		d.sql = cmdtostring(cmd, args...)
		d.mu.Lock()
		d.f.WriteString(fmt.Sprintf("-- %s, reason: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		d.f.WriteString(d.sql + "\n")
		d.f.Sync()
		d.mu.Unlock()
	}
	return 0, err
}

func (d *Db) GetConnections() int {

	return d.conn.Stats().OpenConnections
}

func (d *Db) OpenDebug() *Db {
	d.debug = true
	return d
}

func (d *Db) CloseDebug() *Db {
	d.debug = false
	return d
}

func (d *Db) ping() error {
	return d.conn.Ping()
}

func (d *Db) Update(cmd string, args ...interface{}) (int64, error) {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	err := d.privateTooManyConn()
	if err != nil {
		return 0, err
	}

	result, err := d.conn.ExecContext(d.Ctx, cmd, args...)
	if err != nil {
		return d.execError(err, cmd, args...)
	}

	return result.RowsAffected()
}

func (d *Db) Delete(cmd string, args ...interface{}) (int64, error) {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}

	return d.Update(cmd, args...)
}

func (d *Db) Insert(cmd string, args ...interface{}) (int64, error) {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	err := d.privateTooManyConn()
	if err != nil {
		return 0, err
	}

	result, err := d.conn.ExecContext(d.Ctx, cmd, args...)
	if err != nil {
		return d.execError(err, cmd, args...)
	}

	return result.RowsAffected()
}

func (d *Db) GetSql() string {
	return d.sql
}

func (d *Db) InsertMany(cmd string, args ...interface{}) (int64, error) {
	// sql: insert into test(id, name) values(?,?)  args: interface{}...  1,'t1', 2, 't2', 3, 't3'

	if args == nil {
		return d.Insert(cmd)
	}

	// 先转为为小写
	lowercmd := strings.ToLower(cmd)
	// 找到关键字 values
	tmp_index := strings.Index(lowercmd, " values")
	if tmp_index < 0 {
		return 0, errors.New("insert sql error")
	}
	// 找到关键字 后面的第一个 (
	start_index := strings.Index(cmd[tmp_index:], "(")
	if start_index < 0 {
		return 0, errors.New("sql error: eg: insert into table(name) values(?)")
	}
	end_index := strings.LastIndex(cmd, ")")
	if start_index < 0 {
		return 0, errors.New("sql error: eg: insert into table(name) values(?)")
	}
	value := cmd[tmp_index+start_index : end_index+1]
	//查看一行数据有多少列
	column := 0
	for _, v := range strings.Split(value[1:len(value)-1], ",") {
		opt := strings.Trim(v, " ")
		if opt == "?" {
			column++
		}
	}

	// 总共多少参数
	count := len(args)
	if count%column != 0 {
		return 0, errors.New("args error")
	}
	addcmd := "," + value
	for i := 1; i < count/column; i++ {
		cmd = cmd + addcmd
	}

	return d.Insert(cmd, args...)
}

func (d *Db) GetRows(cmd string, args ...interface{}) (*sql.Rows, error) {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	err := d.privateTooManyConn()
	if err != nil {
		return nil, err
	}

	return d.conn.QueryContext(d.Ctx, cmd, args...)
}

func (d *Db) Close() error {
	//存在并且不为空才关闭
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}

func (d *Db) GetOne(cmd string, args ...interface{}) *Row {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	err := d.privateTooManyConn()
	if err != nil {
		return &Row{err: err}
	}

	return &Row{
		d.conn.QueryRowContext(d.Ctx, cmd, args...), nil}
}

// 还原sql
func cmdtostring(cmd string, args ...interface{}) string {
	cmd = strings.Replace(cmd, "?", "%v", -1)
	if len(args) > 0 {
		newargs := make([]interface{}, 0, len(args))
		for _, v := range args {
			v = fmt.Sprintf("'%v'", v)
			newargs = append(newargs, v)
		}
		return fmt.Sprintf(cmd, newargs...)
	}
	return cmd
}

func (d *Db) privateTooManyConn() error {
	timeout := time.Microsecond * 10
	for d.conn.Stats().OpenConnections >= d.maxConn {
		if timeout.Microseconds() < d.sc.ReadTimeout.Microseconds()/2 {
			time.Sleep(timeout)
			timeout = timeout * 2
		} else {
			return errors.New("read io timeout, more than " + d.sc.ReadTimeout.String())
		}

	}
	return nil
}
