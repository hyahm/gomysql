package gomysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Db struct {
	*sql.DB
	conf      string
	Ctx       context.Context
	sc        *Sqlconfig
	f         *os.File
	mu        *sync.RWMutex
	maxpacket uint64
	maxConn   int
}

func (d *Db) execError(err error, cmd string, args ...interface{}) {

	if d.sc.WriteLogWhenFailed {
		d.mu.Lock()
		d.f.WriteString(fmt.Sprintf("-- %s, reason: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		d.f.WriteString(ToSql(cmd, args...) + ";\n")
		d.f.Sync()
		d.mu.Unlock()
	}
}

func (d *Db) GetConnections() int {

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
	d.Close()
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

func (d *Db) Flush() {
	if d.f != nil {
		d.f.Sync()
		d.f.Close()
	}
}

func (d *Db) Update(cmd string, args ...interface{}) Result {
	res := Result{
		Sql: ToSql(cmd, args...),
	}

	result, err := d.ExecContext(d.Ctx, res.Sql)
	if err != nil {
		res.Err = err
		d.execError(err, res.Sql)
		return res
	}
	res.RowsAffected, res.Err = result.RowsAffected()
	return res
}

func (d *Db) Delete(cmd string, args ...interface{}) Result {
	return d.Update(cmd, args...)
}

func (d *Db) Insert(cmd string, args ...interface{}) Result {
	res := Result{
		Sql: ToSql(cmd, args...),
	}
	result, err := d.ExecContext(d.Ctx, cmd, args...)
	if err != nil {
		res.Err = err
		d.execError(err, res.Sql)
		return res
	}
	res.LastInsertId, res.Err = result.LastInsertId()
	return res
}

func (d *Db) InsertMany(cmd string, args ...interface{}) Result {
	// sql: insert into test(id, name) values(?,?)  args: interface{}...  1,'t1', 2, 't2', 3, 't3'
	// 每次返回的是第一次插入的id
	if args == nil {
		return d.Insert(cmd)
	}
	newcmd, err := formatSql(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.Insert(newcmd, args...)
}

func (d *Db) GetRows(cmd string, args ...interface{}) (*sql.Rows, error) {
	return d.QueryContext(d.Ctx, cmd, args...)
}

func (d *Db) GetOne(cmd string, args ...interface{}) Result {
	result := Result{Sql: ToSql(cmd, args...)}
	row := d.QueryRowContext(d.Ctx, result.Sql)
	if row.Err() != nil {
		result.Err = row.Err()
		return result
	}
	return result
}

func (d *Db) Select(dest interface{}, cmd string, args ...interface{}) Result {
	// db.Select(&value, "select * from test")
	// 传入切片的地址， 根据tag 的 db 自动补充，
	// 最求性能建议还是使用 GetRows or GetOne
	res := Result{Sql: ToSql(cmd, args...)}
	rows, err := d.QueryContext(d.Ctx, res.Sql)
	if err != nil {
		res.Err = err
		return res
	}
	defer rows.Close()
	// 需要设置的值
	res.Err = fill(dest, rows)
	return res
}

func (d *Db) InsertInterfaceWithID(dest interface{}, cmd string, args ...interface{}) Result {
	// $key 和 $value 固定位置固定值
	// db.InsertInterfaceWithID(&value, "insert into test($key)  values($value)")
	res := Result{
		Sql:           ToSql(cmd, args...),
		LastInsertIds: make([]int64, 0),
	}

	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}

	if typ.Kind() == reflect.Struct {
		return d.insertInterface(dest, cmd, args...)
	}
	if typ.Kind() == reflect.Slice {
		// 如果是切片， 那么每个值都做一次处理
		length := value.Len()
		if length == 1 {
			return d.insertInterface(dest, cmd, args...)
		}
		for i := 0; i < length; i++ {
			result := d.insertInterface(value.Index(i).Interface(), cmd, args...)
			res.Sql += ";" + result.Sql
			if result.Err != nil {
				return result
			}
			res.LastInsertIds = append(res.LastInsertIds, result.LastInsertId)
		}
	}
	return res
}

// 插入字段的占位符 $key, $value
func (d *Db) InsertInterfaceWithoutID(dest interface{}, cmd string, args ...interface{}) Result {
	// $key 和 $value 固定位置固定值
	// ID 自增的必须设置 default
	// db.InsertInterfaceWithoutID(&value, "insert into test($key)  values($value)")
	res := Result{}
	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}
	if typ.Kind() == reflect.Struct {
		return d.insertInterface(dest, cmd, args...)
	}

	if typ.Kind() == reflect.Slice {
		// 如果是切片， 那么每个值都做一次处理
		length := value.Len()
		if length == 1 {
			return d.insertInterface(dest, cmd, args...)
		}

		arguments := make([]interface{}, 0)
		for i := 0; i < length; i++ {
			newcmd, newargs, err := insertInterfaceSql(value.Index(i).Interface(), cmd, args...)
			if err != nil {
				res.Err = err
				return res
			}
			cmd = newcmd
			arguments = append(arguments, newargs...)
		}
		return d.InsertMany(cmd, arguments...)
	}
	return res
}

func (d *Db) insertInterface(dest interface{}, cmd string, args ...interface{}) Result {
	// 插入到args之前  dest 是struct或切片的指针
	newcmd, newargs, err := insertInterfaceSql(dest, cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.Insert(newcmd, newargs...)
}

// 还原sql
func ToSql(cmd string, args ...interface{}) string {
	cmd = strings.Replace(cmd, "?", "%v", -1)
	if len(args) > 0 {
		newargs := make([]interface{}, 0, len(args))
		for _, v := range args {
			switch reflect.TypeOf(v).Kind() {
			case reflect.Float32, reflect.Float64, reflect.Bool,
				reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
				reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				v = fmt.Sprintf("%v", v)
			default:
				v = fmt.Sprintf("'%v'", v)
			}

			newargs = append(newargs, v)
		}
		return fmt.Sprintf(cmd, newargs...)
	}
	return cmd
}

// 还原sql
func InToSql(cmd string, args ...interface{}) string {
	cmd, args, err := makeArgs(cmd, args...)
	if err != nil {
		return ""
	}
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

func (d *Db) UpdateInterface(dest interface{}, cmd string, args ...interface{}) Result {
	newcmd, newargs, err := updateInterfaceSql(dest, cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.Update(newcmd, newargs...)
}
