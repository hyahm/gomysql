package gomysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type Tx struct {
	tx    *sql.Tx
	conf  string
	Ctx   context.Context
	sql   string
	debug bool
	sc    *Sqlconfig
	// f       *os.File
	// mu      *sync.RWMutex
	maxConn int
	db      *Db
}

func (d *Db) NewTx(opt *sql.TxOptions) (*Tx, error) {
	tx, err := d.BeginTx(d.Ctx, opt)
	if err != nil {
		return nil, err
	}
	return &Tx{
		tx,
		d.conf,
		d.Ctx,
		d.sql,
		d.debug,
		d.sc,
		d.maxConn,
		d,
	}, nil
}

func (t *Tx) OpenDebug() {
	t.debug = true
}

func (t *Tx) CloseDebug() {
	t.debug = false
}

func (t *Tx) Update(cmd string, args ...interface{}) (int64, error) {
	if t.tx == nil {
		return 0, errors.New("some thing wrong , may be you need close or rallback")
	}
	if t.debug {
		t.sql = cmdtostring(cmd, args...)
	}
	// err := t.privateTooManyConn()
	// if err != nil {
	// 	t.tx = nil
	// 	return 0, err
	// }

	result, err := t.tx.ExecContext(t.Ctx, cmd, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (t *Tx) Delete(cmd string, args ...interface{}) (int64, error) {
	if t.debug {
		t.sql = cmdtostring(cmd, args...)
	}

	return t.Update(cmd, args...)
}

func (t *Tx) Insert(cmd string, args ...interface{}) (int64, error) {
	if t.debug {
		t.sql = cmdtostring(cmd, args...)
	}
	// err := t.privateTooManyConn()
	// if err != nil {
	// 	return 0, err
	// }

	result, err := t.tx.ExecContext(t.Ctx, cmd, args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (t *Tx) GetSql() string {
	return t.sql
}

func (t *Tx) InsertMany(cmd string, args ...interface{}) (int64, error) {
	// sql: insert into test(id, name) values(?,?)  args: interface{}...  1,'t1', 2, 't2', 3, 't3'

	if args == nil {
		return t.Insert(cmd)
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

	return t.Insert(cmd, args...)
}

func (t *Tx) GetRows(cmd string, args ...interface{}) (*sql.Rows, error) {
	if t.debug {
		t.sql = cmdtostring(cmd, args...)
	}
	// err := t.privateTooManyConn()
	// if err != nil {
	// 	return nil, err
	// }

	return t.tx.QueryContext(t.Ctx, cmd, args...)
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) Close() error {
	//存在并且不为空才关闭
	if t != nil {
		return t.Close()
	}
	return nil
}

func (t *Tx) GetOne(cmd string, args ...interface{}) *Row {
	if t.debug {
		t.sql = cmdtostring(cmd, args...)
	}
	// err := t.privateTooManyConn()
	// if err != nil {
	// 	return &Row{err: err}
	// }

	return &Row{
		t.tx.QueryRowContext(t.Ctx, cmd, args...), nil}
}

// func (t *Tx) privateTooManyConn() error {
// 	timeout := time.Microsecond * 10
// 	for t.db.Stats().OpenConnections >= t.maxConn {
// 		if timeout.Microseconds() < t.sc.ReadTimeout.Microseconds()/2 {
// 			time.Sleep(timeout)
// 			timeout = timeout * 2
// 		} else {
// 			return errors.New("read io timeout, more than " + t.sc.ReadTimeout.String())
// 		}

// 	}
// 	return nil
// }
