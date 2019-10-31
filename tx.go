package gomysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type Tx struct {
	tx *sql.Tx
	key string
	db *sql.DB
	conf string
	Ctx context.Context
}

func GetTx(tag string) (*Tx, error) {
	if client == nil {
		return nil, NotInitERROR
	}
	if _, ok := client[tag]; !ok {
		return nil, TAGERROR
	}

	t := &Tx {
		key: tag,
		Ctx: context.Background(),
		conf: client[tag].conf,
		db: client[tag].conn,
	}
	return t,nil
}

func (t *Tx) Begin() (err error) {
	t.tx, err = t.db.Begin()
	if err != nil {
		return
	}
	return
}

func (t *Tx) ping() bool {
	if err := t.db.Ping(); err != nil {
		return false
	}
	return true
}

func (t *Tx) connDb() error {
	// 先要连接db, 然后才连接tx
	db, err := sql.Open("mysql", t.conf)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}
	t.tx, err = db.Begin()
	if err != nil {
		return err
	}

	return nil
}

func (t *Tx)Update(cmd string, args ...interface{}) (int64, error) {

	if !t.ping() {
		// 重连
		if err := t.connDb(); err != nil {
			panic(err)
		}
	}
	result, err :=  t.tx.ExecContext(t.Ctx, cmd, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()

}

func (t *Tx)Insert(cmd string, args ...interface{}) (int64, error) {
	if !t.ping() {
		// 重连
		if err := t.connDb(); err != nil {
			panic(err)
		}
	}
	result, err :=  t.tx.ExecContext(t.Ctx, cmd, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}


func (t *Tx)InsertMany(cmd string, args []interface{}) (int64, error) {

	if !t.ping() {
		// 重连
		if err := t.connDb(); err != nil {
			panic(err)
		}
	}
	if args == nil {
		return t.Insert(cmd, args)
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

	return t.Insert(cmd, args...)
}

func (t *Tx)GetRows(cmd string, args ...interface{}) (*sql.Rows, error) {
	if !t.ping() {
		if err := t.connDb(); err != nil {
			panic(err)
		}
	}
	return t.tx.QueryContext(t.Ctx, cmd, args...)

}

func (t *Tx)Commit() error {
	//存在并且不为空才关闭
	return t.tx.Commit()
}

func (t *Tx)RollBack() error {
	return t.tx.Rollback()
}

func (t *Tx)Close() error {
	//存在并且不为空才关闭
	return t.db.Close()
}

func (t *Tx)GetOne(cmd string, args ...interface{}) *sql.Row {
	if !t.ping() {
		if err := t.connDb(); err != nil {
			panic(err)
		}
	}
	return t.tx.QueryRowContext(t.Ctx, cmd, args...)
}
