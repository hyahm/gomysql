package gomysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

//var conn map[string]*sql.DB
var client map[string]*Db

type Db struct {
	conn *sql.DB
	conf string
	key string
	Ctx context.Context
}

func GetDb(tag string) (*Db, error) {
	if client == nil {
		return nil, NotInitERROR
	}
	if _, ok := client[tag]; !ok {
		return nil, TAGERROR
	}
	return client[tag],nil
}

func (d *Db)connDB() error {
	db, err := sql.Open("mysql", client[d.key].conf)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}
	client[d.key].conn = db
	return nil
}


func (d *Db)GetConnections() int {
	return d.conn.Stats().OpenConnections
}

func (d *Db)ping() bool {
	if err := d.conn.Ping(); err != nil {
		return false
	}
	return true
}


func (d *Db)Update(cmd string, args ...interface{}) (int64, error) {
	if !d.ping() {
		// 重连
		if err := d.connDB(); err != nil {
			panic(err)
		}
	}
	result, err :=  d.conn.ExecContext(d.Ctx, cmd, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *Db)Delete(cmd string, args ...interface{}) (int64, error) {
	//if !d.ping() {
	//	// 重连
	//	if err := d.connDB(); err != nil {
	//		panic(err)
	//	}
	//}
	//result, err :=  d.conn.ExecContext(d.Ctx, cmd, args...)
	//if err != nil {
	//	return 0, err
	//}
	return d.Update(cmd, args)
}


func (d *Db)Insert(cmd string, args ...interface{}) (int64, error) {
	if !d.ping() {
		// 重连
		if err := d.connDB(); err != nil {
			panic(err)
		}
	}
	result, err :=  d.conn.ExecContext(d.Ctx, cmd, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}


func (d *Db)InsertMany(cmd string, args []interface{}) (int64, error) {

	if !d.ping() {
		// 重连
		if err := d.connDB(); err != nil {
			panic(err)
		}
	}
	if args == nil {
		return d.Insert(cmd)
	}

		//找到括号的内容
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
	for _, v := range strings.Split(value, ",") {
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
	if !d.ping() {
		if err := d.connDB(); err != nil {
			panic(err)
		}
	}
	return d.conn.QueryContext(d.Ctx, cmd, args...)

}

func (d *Db)Close() error {
	//存在并且不为空才关闭
	return d.conn.Close()

}

func (d *Db)GetOne(cmd string, args ...interface{}) *sql.Row {
	if !d.ping() {
		if err := d.connDB(); err != nil {
			panic(err)
		}
	}
	return d.conn.QueryRowContext(d.Ctx, cmd, args...)
}

// 还原sql
func cmdtostring(cmd string, args ...interface{}) string {

	var logstr string

	for _, v := range args {
		switch v.(type) {
		case int64:
			logstr = "'" + strconv.FormatInt(v.(int64), 10) + "'"
		case int:
			logstr = "'" + strconv.Itoa(v.(int)) + "'"
		default:
			logstr = "'" + v.(string) + "'"
			//return
		}
		cmd = strings.Replace(cmd, "?", "%s", 1)
		cmd = fmt.Sprintf(cmd, logstr)

	}
	return cmd
}
