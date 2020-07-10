package gomysql

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

var RangeOutErr = errors.New("索引值超出范围")

func makeArgs(cmd string, args ...interface{}) (string, []interface{}, error) {
	vs := make([]interface{}, 0)
	var err error
	for index, value := range args {
		typ := reflect.TypeOf(value)
		if typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice {
			vv := reflect.ValueOf(value)
			cmd, err = replace(cmd, index, vv.Len())
			if err != nil {
				return cmd, vs, err
			}
			for i := 0; i < vv.Len(); i++ {
				vs = append(vs, vv.Index(i).Interface())
			}
		} else {
			vs = append(vs, value)
		}
	}
	return cmd, vs, nil
}

func replace(cmd string, index int, count int) (string, error) {
	m := make([]string, count)
	for j := 0; j < count; j++ {
		m[j] = "?"
	}

	c := strings.Count(cmd, "?")
	if index > c-1 {
		return "", RangeOutErr
	}
	start := 0
	tmp := cmd
	for i := 0; i < c; i++ {
		thisIndex := strings.Index(tmp, "?")
		start = start + thisIndex + 1
		tmp = tmp[thisIndex+1:]
		if i == index {

			return cmd[:start-1] + strings.Join(m, ",") + cmd[start:], nil
		}
	}

	return cmd, nil
}

func (d *Db) UpdateIn(cmd string, args ...interface{}) (int64, error) {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return 0, err
	}
	return d.Update(newcmd, newargs...)
}

func (d *Db) DeleteIn(cmd string, args ...interface{}) (int64, error) {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return 0, err
	}
	return d.Delete(newcmd, newargs...)
}

func (d *Db) GetRowsIn(cmd string, args ...interface{}) (*sql.Rows, error) {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return nil, err
	}
	return d.GetRows(newcmd, newargs...)
}

func (d *Db) GetOneIn(cmd string, args ...interface{}) *sql.Row {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		panic(err)
	}
	return d.GetOne(newcmd, newargs...)
}
