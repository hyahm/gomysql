package gomysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hyahm/golog"
)

type Db struct {
	*sql.DB
	conf      string
	Ctx       context.Context
	sql       string
	debug     bool
	sc        *Sqlconfig
	f         *os.File
	mu        *sync.RWMutex
	maxpacket uint64
	maxConn   int
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

	ch <- struct{}{}
	defer func() {
		<-ch
	}()
	return d.Stats().OpenConnections
}

func (d *Db) Use(dbname string, overWrite ...bool) (*Db, error) {
	// 切换到新的库， 并产生一个新的db
	ow := false
	if len(overWrite) > 0 {
		ow = overWrite[0]
	}
	err := d.CreateDatabase(dbname, ow)
	if err != nil {
		return d, err
	}

	s := &Sqlconfig{
		UserName:                d.sc.UserName,
		Password:                d.sc.Password,
		Host:                    d.sc.Host,
		Port:                    d.sc.Port,
		DbName:                  dbname,
		ClientFoundRows:         d.sc.ClientFoundRows,
		AllowCleartextPasswords: d.sc.AllowCleartextPasswords,
		InterpolateParams:       d.sc.InterpolateParams,
		ColumnsWithAlias:        d.sc.ColumnsWithAlias,
		MultiStatements:         d.sc.MultiStatements,
		ParseTime:               d.sc.ParseTime,
		TLS:                     d.sc.TLS,
		ReadTimeout:             d.sc.ReadTimeout,
		Timeout:                 d.sc.Timeout,
		WriteTimeout:            d.sc.WriteTimeout,
		AllowOldPasswords:       d.sc.AllowOldPasswords,
		Charset:                 d.sc.Charset,
		Loc:                     d.sc.Loc,
		MaxAllowedPacket:        d.sc.MaxAllowedPacket,
		Collation:               d.sc.Collation,
		MaxOpenConns:            d.sc.MaxOpenConns,
		MaxIdleConns:            d.sc.MaxIdleConns,
		ConnMaxLifetime:         d.sc.ConnMaxLifetime,
		WriteLogWhenFailed:      d.sc.WriteLogWhenFailed,
		LogFile:                 d.sc.LogFile,
	}
	return s.conndb(s.GetMysqwlDataSource())
}

func (d *Db) CreateDatabase(dbname string, overWrite bool) error {
	if overWrite {
		d.QueryRow("drop database " + dbname + " if exsits")
	}
	err := d.QueryRow("create database " + dbname).Err()
	if err != nil && err.(*mysql.MySQLError).Number == 1007 {
		return nil
	}
	return err
}

func (d *Db) OpenDebug() {
	d.debug = true
}

func (d *Db) CloseDebug() {
	d.debug = false
}

func (d *Db) Flush() {
	if d.f != nil {
		d.f.Sync()
		d.f.Close()
	}
}

func (d *Db) Update(cmd string, args ...interface{}) (int64, error) {

	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	ch <- struct{}{}
	defer func() {
		<-ch
	}()
	// err := d.privateTooManyConn()
	// if err != nil {
	// 	return 0, err
	// }

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
	ch <- struct{}{}
	defer func() {
		<-ch
	}()
	// err := d.privateTooManyConn()
	// if err != nil {
	// 	return 0, err
	// }
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
	// 每次返回的是第一次插入的id
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
	ch <- struct{}{}
	defer func() {
		<-ch
	}()
	// err := d.privateTooManyConn()
	// if err != nil {
	// 	return nil, err
	// }
	return d.QueryContext(d.Ctx, cmd, args...)
}

func (d *Db) GetOne(cmd string, args ...interface{}) *sql.Row {
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	ch <- struct{}{}
	defer func() {
		<-ch
	}()

	return d.QueryRowContext(d.Ctx, cmd, args...)
}

func (d *Db) Select(dest interface{}, cmd string, args ...interface{}) error {
	// db.Select(&value, "select * from test")
	// 传入切片的地址， 根据tag 的 db 自动补充，
	// 最求性能建议还是使用 GetRows or GetOne
	if d.debug {
		d.sql = cmdtostring(cmd, args...)
	}
	ch <- struct{}{}
	defer func() {
		<-ch
	}()

	rows, err := d.QueryContext(d.Ctx, cmd, args...)
	if err != nil {
		return err
	}
	// 需要设置的值
	value := reflect.ValueOf(dest)
	typ := reflect.TypeOf(dest)
	// cols := 0
	// // json.Unmarshal returns errors for these
	if typ.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
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

					case reflect.Slice, reflect.Struct:
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
					golog.Info("can not set: ", index)
				}
			}

		}
		if !isArr {
			if isPtr {
				value.Elem().Elem().Set(new)
			} else {
				value.Elem().Set(new)
			}

			return nil
		} else {
			if isPtr {
				ss = reflect.Append(ss, new.Addr())
			} else {
				ss = reflect.Append(ss, new)

			}
		}
	}
	value.Elem().Set(ss)
	return nil
}

func (d *Db) InsertInterfaceWithID(dest interface{}, cmd string, args ...interface{}) ([]int64, error) {
	// $key 和 $value 固定位置固定值
	// db.InsertInterfaceWithID(&value, "insert into test($key)  values($value)")
	if !strings.Contains(cmd, "$key") {
		return nil, errors.New("not found placeholders $key")
	}

	if !strings.Contains(cmd, "$value") {
		return nil, errors.New("not found placeholders $value")
	}
	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}

	if typ.Kind() == reflect.Struct {
		id, err := d.insertInterface(dest, cmd, args...)
		return []int64{id}, err
	}
	ids := make([]int64, 0)
	if typ.Kind() == reflect.Slice {
		// 如果是切片， 那么每个值都做一次处理
		length := value.Len()

		for i := 0; i < length; i++ {
			id, err := d.insertInterface(value.Index(i).Interface(), cmd, args...)
			if err != nil {
				return ids, err
			}
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// 插入字段的占位符 $key, $value
func (d *Db) InsertInterfaceWithoutID(dest interface{}, cmd string, args ...interface{}) error {
	// $key 和 $value 固定位置固定值
	// db.InsertInterfaceWithoutID(&value, "insert into test($key)  values($value)")
	if !strings.Contains(cmd, "$key") {
		return errors.New("not found placeholders $key")
	}

	if !strings.Contains(cmd, "$value") {
		return errors.New("not found placeholders $value")
	}
	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}
	if typ.Kind() == reflect.Struct {
		_, err := d.insertInterface(dest, cmd, args...)
		return err
	}

	if typ.Kind() == reflect.Slice {
		// 如果是切片， 那么每个值都做一次处理
		length := value.Len()

		for i := 0; i < length; i++ {
			_, err := d.insertInterface(value.Index(i).Interface(), cmd, args...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Db) insertInterface(dest interface{}, cmd string, args ...interface{}) (int64, error) {
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
				if value.Field(i) == reflect.ValueOf("") && strings.Contains(key, "omitempty") {
					continue
				}
				keys = append(keys, signs[0])
				placeholders = append(placeholders, "?")
				values = append(values, value.Field(i).Interface())
			case reflect.Int64,
				reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Float32, reflect.Float64:
				if value.Field(i) == reflect.ValueOf(0) && strings.Contains(key, "omitempty") {
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
					values = append(values, "")
				} else {
					if value.Field(i).Len() == 0 && !strings.Contains(key, "omitempty") {
						continue
					}
					keys = append(keys, signs[0])
					placeholders = append(placeholders, "?")
					send, err := json.Marshal(value.Field(i).Interface())
					if err != nil {
						values = append(values, "")
						continue
					}
					values = append(values, send)
				}
			case reflect.Ptr:
				if value.Field(i).IsNil() {
					if !strings.Contains(key, "omitempty") {
						continue
					}
					keys = append(keys, signs[0])
					placeholders = append(placeholders, "?")
					values = append(values, "")
				} else {
					keys = append(keys, signs[0])
					placeholders = append(placeholders, "?")
					send, err := json.Marshal(value.Field(i).Interface())
					if err != nil {
						values = append(values, "")
						continue
					}
					values = append(values, send)
				}
			case reflect.Struct:
				keys = append(keys, signs[0])
				placeholders = append(placeholders, "?")
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
	}

	cmd = strings.Replace(cmd, "$key", strings.Join(keys, ","), 1)
	cmd = strings.Replace(cmd, "$value", strings.Join(placeholders, ","), 1)
	newargs := append(values, args...)
	return d.Insert(cmd, newargs...)
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

func (d *Db) UpdateInterface(dest interface{}, cmd string, args ...interface{}) (int64, error) {
	// $set 固定位置固定值
	// db.UpdateInterface(&value, "update test set $set where id=1")
	// 插入到args之前  dest 是struct或切片的指针
	if !strings.Contains(cmd, "$set") {
		return 0, errors.New("not found placeholders $set")
	}

	// ？号
	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)

	if typ.Kind() == reflect.Ptr {
		value = value.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return 0, errors.New("dest must ptr or struct ")
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
		case reflect.Int64,
			reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Float32, reflect.Float64:
			if value.Field(i) == reflect.ValueOf(0) && !strings.Contains(key, "force") {
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
	return d.Update(cmd, newargs...)
}
