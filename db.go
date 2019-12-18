package gomysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type Db struct {
	conn *sql.DB
	conf string
	Ctx context.Context
	sql string
	debug bool
}

func (d *Db) conndb() (*Db, error) {
	conn,err := sql.Open("mysql", d.conf)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	d.conn = conn
	return d, nil
}

func (d *Db)GetConnections() int {
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

func (d *Db)ping() error {
	return d.conn.Ping()
}

func (d *Db)Update(cmd string, args ...interface{}) (int64, error) {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	if err := d.ping(); err != nil  {
		// 重连
		if d, err = d.conndb();err != nil {
			return 0,err
		}
	}
	result, err :=  d.conn.ExecContext(d.Ctx, cmd, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *Db)Delete(cmd string, args ...interface{}) (int64, error) {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	return d.Update(cmd, args...)
}


func (d *Db)Insert(cmd string, args ...interface{}) (int64, error) {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	if err := d.ping(); err != nil  {
		// 重连
		if d, err = d.conndb();err != nil {
			return 0,err
		}
	}
	result, err :=  d.conn.ExecContext(d.Ctx, cmd, args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
	//return 0, nil
}

func (d *Db) GetSql() string {
	return d.sql
}


func (d *Db)InsertMany(cmd string, args ...interface{}) (int64, error) {

	if err := d.ping(); err != nil  {
		// 重连
		if d, err = d.conndb();err != nil {
			return 0,err
		}
	}
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
	value := cmd[tmp_index+start_index: end_index+1]
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
	if count % column != 0 {
		return 0, errors.New("args error")
	}
	addcmd := "," + value
	for i := 1; i < count / column; i++ {
		cmd = cmd + addcmd
	}

	return d.Insert(cmd, args...)
}

func (d *Db)GetRows(cmd string, args ...interface{}) (*sql.Rows, error) {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	if err := d.ping(); err != nil  {
		// 重连
		if d, err = d.conndb();err != nil {
			return nil,err
		}
	}
	return d.conn.QueryContext(d.Ctx, cmd, args...)

}

func (d *Db)Close() error {
	//存在并且不为空才关闭
	if d.conn != nil {
		return d.conn.Close()
	}
	return  nil
}

func (d *Db)GetOne(cmd string, args ...interface{}) (*sql.Row,error) {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	if err := d.ping(); err != nil  {
		// 重连
		if d, err = d.conndb();err != nil {
			return nil, err
		}
	}
	return d.conn.QueryRowContext(d.Ctx, cmd, args...), nil
}

// 还原sql
func cmdtostring(cmd string, args ...interface{}) string {
	cmd = strings.Replace(cmd, "?", "%v", -1)
	for _, v := range args {
		v = fmt.Sprintf("'%v'", v)
	}
	return fmt.Sprintf(cmd, args...)
}
