package gomysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	*sql.DB
	conf    string
	Ctx     context.Context
	sql     string
	debug   bool
	sc      *Sqlconfig
	f       *os.File
	mu      *sync.RWMutex
	maxConn int
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

	return d.Stats().OpenConnections
}

func (d *Db) OpenDebug() {
	d.debug = true
}

func (d *Db) CloseDebug() {
	d.debug = false
}

func (d *Db) ping() error {
	return d.Ping()
}

func (d *Db) Update(cmd string, args ...interface{}) (int64, error) {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	err := d.privateTooManyConn()
	if err != nil {
		return 0, err
	}

	result, err := d.ExecContext(d.Ctx, cmd, args...)
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

	result, err := d.ExecContext(d.Ctx, cmd, args...)
	if err != nil {
		return d.execError(err, cmd, args...)
	}

	return result.LastInsertId()
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

	return d.QueryContext(d.Ctx, cmd, args...)
}

func (d *Db) Close() error {
	//存在并且不为空才关闭
	if d != nil {
		return d.Close()
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
		d.QueryRowContext(d.Ctx, cmd, args...), nil}
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
	for d.Stats().OpenConnections >= d.maxConn {
		if timeout.Microseconds() < d.sc.ReadTimeout.Microseconds()/2 {
			time.Sleep(timeout)
			timeout = timeout * 2
		} else {
			return errors.New("read io timeout, more than " + d.sc.ReadTimeout.String())
		}

	}
	return nil
}
