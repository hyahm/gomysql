package gomysql

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"reflect"
	"strings"
)

type Tx struct {
	tx      *sql.Tx
	conf    string
	Ctx     context.Context
	sc      *Sqlconfig
	maxConn int
	db      *Db
	debug   bool
}

func (d *Db) NewTx() (*Tx, error) {
	tx, err := d.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{
		tx,
		d.conf,
		d.Ctx,
		d.sc,
		d.maxConn,
		d,
		d.debug,
	}, nil
}

func (t *Tx) Update(cmd string, args ...interface{}) Result {
	res := Result{}
	if t.debug {
		res.Sql = ToSql(cmd, args...)
	}
	if t.tx == nil {
		res.Err = errors.New("some thing wrong , may be you need close or rallback")
		return res
	}

	result, err := t.tx.ExecContext(t.Ctx, cmd, args...)
	if err != nil {
		res.Err = err
		return res
	}
	affect, _ := result.RowsAffected()
	res.LastInsertId = affect
	return res
}

func (t *Tx) Delete(cmd string, args ...interface{}) Result {
	return t.Update(cmd, args...)
}

func (t *Tx) Insert(cmd string, args ...interface{}) Result {
	res := Result{}
	if t.debug {
		res.Sql = ToSql(cmd, args...)
	}
	result, err := t.tx.ExecContext(t.Ctx, cmd, args...)
	if err != nil {
		res.Err = err
		return res
	}
	res.LastInsertId, res.Err = result.LastInsertId()
	return res
}

func (t *Tx) InsertMany(cmd string, args ...interface{}) Result {
	// sql: insert into test(id, name) values(?,?)  args: interface{}...  1,'t1', 2, 't2', 3, 't3'
	if args == nil {
		return t.Insert(cmd)
	}
	newcmd, err := formatSql(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return t.Insert(newcmd, args...)
}

func (t *Tx) GetRows(cmd string, args ...interface{}) (*sql.Rows, error) {
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

func (t *Tx) GetOne(cmd string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(t.Ctx, cmd, args...)
}

func (t *Tx) UpdateIn(cmd string, args ...interface{}) Result {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return t.Update(newcmd, newargs...)
}

func (t *Tx) InsertIn(cmd string, args ...interface{}) Result {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return t.Insert(newcmd, newargs...)
}

func (t *Tx) DeleteIn(cmd string, args ...interface{}) Result {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return t.Delete(newcmd, newargs...)
}

func (t *Tx) GetRowsIn(cmd string, args ...interface{}) (*sql.Rows, error) {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return nil, err
	}
	return t.GetRows(newcmd, newargs...)
}

func (t *Tx) GetOneIn(cmd string, args ...interface{}) *sql.Row {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		log.Println(err)
		return nil
	}
	return t.GetOne(newcmd, newargs...)
}

func (t *Tx) UpdateInterfaceIn(dest interface{}, cmd string, args ...interface{}) Result {
	// $set 固定位置固定值
	// db.UpdateInterfaceIn(&value, "update test set $set where id in (?)", [])
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return t.UpdateInterface(dest, newcmd, newargs...)
}

func (t *Tx) InsertInterfaceWithoutIDIn(dest interface{}, cmd string, args ...interface{}) Result {
	// $key 和 $value 固定位置固定值
	// db.InsertInterfaceWithoutIDIn(&value, "insert into test($key)  values($value)  ", [])

	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}

	return t.InsertInterfaceWithoutID(dest, newcmd, newargs...)
}

func (t *Tx) InsertInterfaceWithIDIn(dest interface{}, cmd string, args ...interface{}) Result {
	// $key 和 $value 固定位置固定值
	// db.InsertInterfaceWithIDIn(&value, "insert into test($key)  values($value) where a in (?)" [])
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return t.InsertInterfaceWithID(dest, newcmd, newargs...)
}

func (t *Tx) SelectIn(dest interface{}, cmd string, args ...interface{}) Result {
	// db.SelectIn(&value, "select * from test  where a in (?)", [])
	// 传入切片的地址， 根据tag 的 db 自动补充，
	// 最求性能建议还是使用 GetRows or GetOne
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return t.Select(dest, newcmd, newargs...)
}

func (t *Tx) Select(dest interface{}, cmd string, args ...interface{}) Result {
	// db.Select(&value, "select * from test")
	// 传入切片的地址， 根据tag 的 db 自动补充，
	// 最求性能建议还是使用 GetRows or GetOne
	res := Result{}
	if t.debug {
		res.Sql = ToSql(cmd, args...)
	}
	rows, err := t.tx.QueryContext(t.Ctx, cmd, args...)
	if err != nil {
		res.Err = err
		return res
	}
	defer rows.Close()
	// 需要设置的值
	res.Err = fill(dest, rows)
	return res
}

func (t *Tx) InsertInterfaceWithID(dest interface{}, cmd string, args ...interface{}) Result {
	// $key 和 $value 固定位置固定值
	// db.InsertInterfaceWithID(&value, "insert into test($key)  values($value)")
	res := Result{
		LastInsertIds: make([]int64, 0),
	}
	if !strings.Contains(cmd, "$key") {
		res.Err = errors.New("not found placeholders $key")
		return res
	}

	if !strings.Contains(cmd, "$value") {
		res.Err = errors.New("not found placeholders $value")
		return res
	}
	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}

	if typ.Kind() == reflect.Struct {
		return t.insertInterface(dest, cmd, args...)
	}
	if typ.Kind() == reflect.Slice {
		// 如果是切片， 那么每个值都做一次处理
		length := value.Len()
		if length == 1 {
			return t.insertInterface(dest, cmd, args...)
		}
		for i := 0; i < length; i++ {
			result := t.insertInterface(value.Index(i).Interface(), cmd, args...)
			res.Sql += ";" + result.Sql
			if result.Err != nil {
				return res
			}
			res.LastInsertIds = append(res.LastInsertIds, result.LastInsertId)
		}
	} else {
		res.Err = ErrNotSupport
	}
	return res
}

// 插入字段的占位符 $key, $value
func (t *Tx) InsertInterfaceWithoutID(dest interface{}, cmd string, args ...interface{}) Result {
	// $key 和 $value 固定位置固定值
	// ID 自增的必须设置 default
	// db.InsertInterfaceWithoutID(&value, "insert into test($key)  values($value)")
	if !strings.Contains(cmd, "$key") {
		return Result{Err: errors.New("not found placeholders $key")}
	}

	if !strings.Contains(cmd, "$value") {
		return Result{Err: errors.New("not found placeholders $value")}
	}
	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}
	if typ.Kind() == reflect.Struct {
		return t.insertInterface(dest, cmd, args...)
	}

	if typ.Kind() == reflect.Slice {
		// 如果是切片， 那么每个值都做一次处理
		length := value.Len()
		if length == 1 {
			return t.insertInterface(dest, cmd, args...)
		}
		arguments := make([]interface{}, 0)
		for i := 0; i < length; i++ {
			newcmd, newargs, err := insertInterfaceSql(value.Index(i).Interface(), cmd, args...)
			if err != nil {
				return Result{Err: err}
			}
			cmd = newcmd
			arguments = append(arguments, newargs...)
		}
		return t.InsertMany(cmd, arguments...)
	} else {
		return Result{Err: ErrNotSupport}
	}
}

func (t *Tx) insertInterface(dest interface{}, cmd string, args ...interface{}) Result {
	newcmd, newargs, err := insertInterfaceSql(dest, cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return t.Insert(newcmd, newargs...)
}

func (t *Tx) UpdateInterface(dest interface{}, cmd string, args ...interface{}) Result {
	// $set 固定位置固定值
	// db.UpdateInterface(&value, "update test set $set where id=1")
	// 插入到args之前  dest 是struct或切片的指针
	newcmd, newargs, err := updateInterfaceSql(dest, cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return t.Update(newcmd, newargs...)
}
