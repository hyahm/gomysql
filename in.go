package gomysql

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/hyahm/golog"
)

// RangeOutErr 错误信息
var RangeOutErr = errors.New("索引值超出范围")

// IgnoreWords 如果删除in的字段， 前面如果出现下面的关键字， 也删除
var IgnoreWords = []string{"or", "and", "where", "on"}

func makeArgs(cmd string, args ...interface{}) (string, []interface{}, error) {
	vs := make([]interface{}, 0)
	var err error
	for index, value := range args {
		typ := reflect.TypeOf(value)
		vv := reflect.ValueOf(value)
		if typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice {
			golog.Info(vv.Len())
			l := vv.Len()
			if l == 0 {
				// 删除此条件
				cmd, err = findStrIndex(cmd, index, true)
				if err != nil {
					return cmd, vs, err
				}
			} else if l == 1 {
				cmd, err = findStrIndex(cmd, index, false)
				if err != nil {
					return cmd, vs, err
				}
			} else {
				for i := 0; i < l; i++ {
					vs = append(vs, vv.Index(i).Interface())
				}
				cmd, err = replace(cmd, index, l)
				if err != nil {
					return cmd, vs, err
				}
			}

		} else {
			// 不是数组的话， 直接返回
			vs = append(vs, value)
		}
	}
	return cmd, vs, nil
}

func findStrIndex(cmd string, pos int, del bool) (string, error) {
	count := strings.Count(cmd, "?")
	tmp := cmd
	lastcmd := ""
	start := 0
	for i := 0; i < count; i++ {
		thisIndex := strings.Index(tmp, "?")
		start = start + thisIndex + 1
		tmp = tmp[thisIndex+1:]
		if i == pos {
			// 找到前面的(
			ksindex := strings.LastIndex(cmd[:start], "(")
			klindex := strings.LastIndex(cmd[start:], ")")
			if ksindex <= 0 {
				return "", errors.New("sql error")
			}
			if del {
				inIndex := strings.LastIndex(cmd[:start], "in")
				// 去掉空格
				aa := strings.Trim(cmd[:inIndex], " ")
				// 替换成=
				// 找到前面一个空格位置
				spaceIndex := strings.LastIndex(aa, " ")
				// 去掉空格 继续找前面的
				bb := strings.Trim(cmd[:spaceIndex], " ")
				tIndex := strings.LastIndex(bb, " ")
				lastStr := strings.Trim(cmd[tIndex:spaceIndex], " ")
				golog.Infof("-%s--", lastStr)
				// 再次查找前面的， 如果是or 或者 and ，wher, on
				for _, word := range IgnoreWords {
					if word == lastStr {
						lastcmd = cmd[:tIndex] + tmp[1:]
						goto endloop
					}
				}

				lastcmd = cmd[:spaceIndex] + tmp[1:]

			} else {
				inIndex := strings.LastIndex(cmd[:start], "in")

				for j := 0; j < len(cmd); j++ {
					if inIndex == j {
						lastcmd += "="
						continue
					}
					if ksindex == j || j == klindex+start || inIndex+1 == j {
						continue
					}
					lastcmd += cmd[j : j+1]
				}
			}
			// 1, 去掉括号， 修改in为 =

			break

		}
	}
endloop:
	return lastcmd, nil
}

func replace(cmd string, index int, count int) (string, error) {
	// 替换？,
	golog.Info("index: ", index)
	// index: 第几个问号开始替换
	// count: 替换多少次
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
