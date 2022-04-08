package gomysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type Tx struct {
	tx      *sql.Tx
	conf    string
	Ctx     context.Context
	sc      *Sqlconfig
	maxConn int
	db      *Db
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
	}, nil
}

func (t *Tx) Update(cmd string, args ...interface{}) Result {
	res := Result{
		Sql: ToSql(cmd, args...),
	}
	if t.tx == nil {
		res.Err = errors.New("some thing wrong , may be you need close or rallback")
		return res
	}

	result, err := t.tx.ExecContext(t.Ctx, res.Sql)
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
	res := Result{
		Sql: ToSql(cmd, args...),
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
	res := Result{
		Sql: ToSql(cmd, args...),
	}
	rows, err := t.tx.QueryContext(t.Ctx, cmd, args...)
	if err != nil {
		res.Err = err
		return res
	}
	defer rows.Close()
	// 需要设置的值
	value := reflect.ValueOf(dest)
	typ := reflect.TypeOf(dest)
	// cols := 0
	// // json.Unmarshal returns errors for these
	if typ.Kind() != reflect.Ptr {
		res.Err = errors.New("must pass a pointer, not a value, to StructScan destination")
		return res
	}
	// stt 是数组基础数据结构

	typ = typ.Elem()
	// 判断是否是数组
	isArr := false
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
		isArr = true
	}
	// 标识最后的接受体是指针还是结构体
	isPtr := false
	if typ.Kind() == reflect.Ptr {
		isPtr = true
		typ = typ.Elem()
	}
	// ss 是切片
	ss := value.Elem()
	names := make(map[string]int)
	cls, _ := rows.Columns()
	for i, v := range cls {
		names[v] = i
	}

	vals := make([][]byte, len(cls))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cls))
	//这里scans引用vals，把数据填充到[]byte里
	for k := range vals {
		scans[k] = &vals[k]
	}
	for rows.Next() {
		// scan into the struct field pointers and append to our results
		err = rows.Scan(scans...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		new := reflect.New(typ)
		if !isPtr {
			new = new.Elem()
		}
		if new.Type().Kind() == reflect.Ptr {
			new = new.Elem()
		}
		for index := 0; index < typ.NumField(); index++ {
			dbname := typ.Field(index).Tag.Get("db")
			tags := strings.Split(dbname, ",")
			if len(tags) < 0 {
				continue
			}
			if tags[0] == "" {
				continue
			}

			if v, ok := names[tags[0]]; ok {
				if new.Field(index).CanSet() {
					// 判断这一列的值
					kind := new.Field(index).Kind()
					b := *(scans[v]).(*[]byte)
					switch kind {
					case reflect.String:

						new.Field(index).SetString(string(b))
					case reflect.Int64:
						i64, _ := strconv.ParseInt(string(b), 10, 64)
						new.Field(index).SetInt(i64)
					case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32:
						i, _ := strconv.Atoi(string(b))
						new.Field(index).Set(reflect.ValueOf(i))

					case reflect.Bool:
						t, _ := strconv.ParseBool(string(b))
						new.Field(index).SetBool(t)

					case reflect.Float32:
						f64, _ := strconv.ParseFloat(string(b), 32)
						new.Field(index).SetFloat(f64)

					case reflect.Float64:
						f64, _ := strconv.ParseFloat(string(b), 64)
						new.Field(index).SetFloat(f64)

					case reflect.Struct, reflect.Slice:
						j := reflect.New(new.Field(index).Type())
						json.Unmarshal(b, j.Interface())
						new.Field(index).Set(j.Elem())

					case reflect.Ptr:
						j := reflect.New(new.Field(index).Type())
						json.Unmarshal(b, j.Interface())
						new.Field(index).Set(j)
					default:
						fmt.Println("not support , you can add issue: ", kind)
					}
				} else {
					fmt.Println("can not set: ", index)
				}
			}

		}
		if !isArr {
			if isPtr {
				value.Elem().Elem().Set(new)
			} else {
				value.Elem().Set(new)
			}

			return res
		} else {
			if isPtr {
				ss = reflect.Append(ss, new.Addr())
			} else {
				ss = reflect.Append(ss, new)

			}
		}
	}
	value.Elem().Set(ss)
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

		for i := 0; i < length; i++ {
			result := t.insertInterface(value.Index(i).Interface(), cmd, args...)
			if result.Err != nil {
				return res
			}
			res.LastInsertIds = append(res.LastInsertIds, result.LastInsertId)
		}
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

		for i := 0; i < length; i++ {
			result := t.insertInterface(value.Index(i).Interface(), cmd, args...)
			if result.Err != nil {
				return result
			}
		}
	}
	return Result{Err: nil}
}

func (t *Tx) insertInterface(dest interface{}, cmd string, args ...interface{}) Result {
	// 插入到args之前  dest 是struct或切片的指针
	values := make([]interface{}, 0)
	keys := make([]string, 0)
	// ？号
	placeholders := make([]string, 0)
	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)

	if typ.Kind() == reflect.Ptr {
		value = value.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Struct {
		// 如果是struct， 执行插入
		for i := 0; i < value.NumField(); i++ {
			key := typ.Field(i).Tag.Get("db")
			if key == "" {
				continue

			}
			signs := strings.Split(key, ",")
			kind := value.Field(i).Kind()
			switch kind {
			case reflect.String:
				if value.Field(i) == reflect.ValueOf("") && strings.Contains(key, "default") {
					continue
				}
				keys = append(keys, signs[0])
				placeholders = append(placeholders, "?")
				values = append(values, value.Field(i).Interface())
			case reflect.Int64,
				reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Float32, reflect.Float64:
				if value.Field(i) == reflect.ValueOf(0) && strings.Contains(key, "default") {
					continue
				}

				keys = append(keys, signs[0])
				placeholders = append(placeholders, "?")
				values = append(values, value.Field(i).Interface())
			case reflect.Bool:
				keys = append(keys, signs[0])
				placeholders = append(placeholders, "?")
				values = append(values, value.Field(i).Interface())
			case reflect.Slice:
				if value.Field(i).IsNil() {
					keys = append(keys, signs[0])
					placeholders = append(placeholders, "?")
					values = append(values, "[]")
				} else {
					if value.Field(i).Len() == 0 && !strings.Contains(key, "default") {
						continue
					}

					keys = append(keys, signs[0])
					placeholders = append(placeholders, "?")
					send, err := json.Marshal(value.Field(i).Interface())
					if err != nil {
						values = append(values, "[]")
						continue
					}
					values = append(values, send)
				}
			case reflect.Ptr:
				if value.Field(i).IsNil() {
					if !strings.Contains(key, "default") {
						continue
					}
					keys = append(keys, signs[0])
					placeholders = append(placeholders, "?")
					values = append(values, "{}")
				} else {
					keys = append(keys, signs[0])
					placeholders = append(placeholders, "?")
					send, err := json.Marshal(value.Field(i).Interface())
					if err != nil {
						values = append(values, "{}")
						continue
					}
					values = append(values, send)
				}
			case reflect.Struct:
				keys = append(keys, signs[0])
				placeholders = append(placeholders, "?")
				send, err := json.Marshal(value.Field(i).Interface())
				if err != nil {
					values = append(values, "{}")
					continue
				}
				values = append(values, send)
			default:
				fmt.Println("not support , you can add issue: ", kind)
			}
		}
	}

	cmd = strings.Replace(cmd, "$key", strings.Join(keys, ","), 1)
	cmd = strings.Replace(cmd, "$value", strings.Join(placeholders, ","), 1)
	newargs := append(values, args...)
	return t.Insert(cmd, newargs...)
}

func (t *Tx) UpdateInterface(dest interface{}, cmd string, args ...interface{}) Result {
	// $set 固定位置固定值
	// db.UpdateInterface(&value, "update test set $set where id=1")
	// 插入到args之前  dest 是struct或切片的指针
	if !strings.Contains(cmd, "$set") {
		return Result{Err: errors.New("not found placeholders $set")}
	}

	// ？号
	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)

	if typ.Kind() == reflect.Ptr {
		value = value.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return Result{Err: errors.New("dest must ptr or struct ")}
	}
	values := make([]interface{}, 0)
	keys := make([]string, 0)
	// 如果是struct， 执行插入
	for i := 0; i < value.NumField(); i++ {
		key := typ.Field(i).Tag.Get("db")
		if key == "" {
			continue
		}
		signs := strings.Split(key, ",")
		kind := value.Field(i).Kind()
		switch kind {
		case reflect.String:
			if value.Field(i).Interface().(string) == "" && !strings.Contains(key, "force") {
				continue
			}
			keys = append(keys, signs[0]+"=?")
			values = append(values, value.Field(i).Interface())
		case reflect.Int64, reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32:
			if value.Field(i).Int() == 0 && !strings.Contains(key, "force") {
				continue
			}
			keys = append(keys, signs[0]+"=?")
			values = append(values, value.Field(i).Interface())
		case reflect.Float32, reflect.Float64:
			if value.Field(i).Float() == 0 && !strings.Contains(key, "force") {
				continue
			}
			keys = append(keys, signs[0]+"=?")
			values = append(values, value.Field(i).Interface())
		case reflect.Uint64, reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32:
			if value.Field(i).Uint() == 0 && !strings.Contains(key, "force") {
				continue
			}
			keys = append(keys, signs[0]+"=?")
			values = append(values, value.Field(i).Interface())
		case reflect.Bool:
			keys = append(keys, signs[0]+"=?")
			values = append(values, value.Field(i).Interface())
		case reflect.Slice:
			if value.Field(i).IsNil() {
				if !strings.Contains(key, "force") {
					continue
				}
				keys = append(keys, signs[0]+"=?")
				values = append(values, "")
			} else {
				if value.Field(i).Len() == 0 && !strings.Contains(key, "force") {
					continue
				}
				keys = append(keys, signs[0]+"=?")
				send, err := json.Marshal(value.Field(i).Interface())
				if err != nil {
					values = append(values, "")
					continue
				}
				values = append(values, send)
			}
		case reflect.Ptr:
			if value.Field(i).IsNil() {
				if !strings.Contains(key, "force") {
					continue
				}
				keys = append(keys, signs[0]+"=?")
				values = append(values, "")
			} else {
				keys = append(keys, signs[0]+"=?")
				send, err := json.Marshal(value.Field(i).Interface())
				if err != nil {
					values = append(values, "")
					continue
				}
				values = append(values, send)
			}
		case reflect.Struct:
			keys = append(keys, signs[0]+"=?")
			send, err := json.Marshal(value.Field(i).Interface())
			if err != nil {
				values = append(values, "")
				continue
			}
			values = append(values, send)
		default:
			fmt.Println("not support , you can add issue: ", kind)
		}
	}

	cmd = strings.Replace(cmd, "$set", strings.Join(keys, ","), 1)
	newargs := append(values, args...)
	return t.Update(cmd, newargs...)
}
