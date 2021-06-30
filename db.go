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

	d.sc.DbName = dbname
	return d.sc.NewDb()
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

func (d *Db) Close() error {
	//存在并且不为空才关闭
	defer func() {
		for {
			<-ch
		}
	}()
	if d != nil {
		return d.Close()
	}
	return nil
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
		golog.Error(err)
		return err
	}
	// 需要设置的值
	value := reflect.ValueOf(dest)
	// cols := 0
	// // json.Unmarshal returns errors for these
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	// stt 是数组基础数据结构

	stt := value.Type().Elem()

	if stt.Kind() == reflect.Slice {
		stt = stt.Elem()
	}
	isPtr := false

	if stt.Kind() == reflect.Ptr {
		isPtr = true
		stt = stt.Elem()
	}

	aa := value.Elem()
	names := make(map[string]int)
	cls, err := rows.Columns()
	if err != nil {
		golog.Error(err)
	}

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
		new := reflect.New(stt)
		if !isPtr {
			new = new.Elem()
		}
		newvalue := new.Elem()
		for index := 0; index < stt.NumField(); index++ {
			dbname := stt.Field(index).Tag.Get("db")
			if dbname == "" {
				continue
			}

			if v, ok := names[dbname]; ok {
				if newvalue.Field(index).CanSet() {
					// 判断这一列的值
					kind := newvalue.Field(index).Kind()
					b := *(scans[v]).(*[]byte)
					switch kind {
					case reflect.String:
						newvalue.Field(index).SetString(string(b))
					case reflect.Int64:
						i64, _ := strconv.ParseInt(string(b), 10, 64)
						newvalue.Field(index).SetInt(i64)
					case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32:
						i, _ := strconv.Atoi(string(b))
						newvalue.Field(index).Set(reflect.ValueOf(i))

					case reflect.Bool:
						t, _ := strconv.ParseBool(string(b))
						newvalue.Field(index).SetBool(t)

					case reflect.Float32:
						f64, _ := strconv.ParseFloat(string(b), 32)
						newvalue.Field(index).SetFloat(f64)

					case reflect.Float64:
						f64, _ := strconv.ParseFloat(string(b), 64)
						newvalue.Field(index).SetFloat(f64)

					case reflect.Slice, reflect.Struct:
						j := reflect.New(newvalue.Field(index).Type())
						golog.Info(j.Type().Kind())
						err = json.Unmarshal(b, j.Interface())
						if err != nil {
							golog.Error(err)
						}
						newvalue.Field(index).Set(j.Elem())

					case reflect.Ptr:
						j := reflect.New(newvalue.Field(index).Type())
						golog.Info(j.Type().Kind())
						err = json.Unmarshal(b, j.Interface())
						if err != nil {
							golog.Error(err)
						}
						newvalue.Field(index).Set(j)
					default:
						golog.Info(kind)
					}
					// new.Field(i).Set()
				}
			}

		}
		golog.Info(new)
		aa = reflect.Append(aa, new)
	}
	value.Elem().Set(aa)
	return nil
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
