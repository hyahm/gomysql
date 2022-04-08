package gomysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// RangeOutErr 错误信息
var ErrRangeOut = errors.New("索引值超出范围")

// IgnoreWords 如果删除in的字段， 前面如果出现下面的关键字， 也删除
var IgnoreWords = []string{"or", "and"}

// 如果后面有or 或者 and，只用删除or 或者and 那么删除前面的where， on 关键字
var IgnoreCondition = []string{"where", "on"}

func (d *Db) UpdateIn(cmd string, args ...interface{}) Result {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.Update(newcmd, newargs)
}

func (d *Db) InsertIn(cmd string, args ...interface{}) Result {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.Insert(newcmd, newargs)
}

func (d *Db) DeleteIn(cmd string, args ...interface{}) Result {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.Delete(newcmd, newargs)
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
		log.Println(err)
		return nil
	}
	return d.GetOne(newcmd, newargs...)
}

func (d *Db) UpdateInterfaceIn(dest interface{}, cmd string, args ...interface{}) Result {
	// $set 固定位置固定值
	// db.UpdateInterfaceIn(&value, "update test set $set where id in (?)", [])
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.UpdateInterface(dest, newcmd, newargs...)
}

func (d *Db) InsertInterfaceWithoutIDIn(dest interface{}, cmd string, args ...interface{}) Result {
	// $key 和 $value 固定位置固定值
	// db.InsertInterfaceWithoutIDIn(&value, "insert into test($key)  values($value)  ", [])
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.InsertInterfaceWithoutID(dest, newcmd, newargs...)
}

func (d *Db) InsertInterfaceWithIDIn(dest interface{}, cmd string, args ...interface{}) Result {
	// $key 和 $value 固定位置固定值
	// db.InsertInterfaceWithIDIn(&value, "insert into test($key)  values($value) where a in (?)" [])
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.InsertInterfaceWithID(dest, newcmd, newargs...)
}

func (d *Db) SelectIn(dest interface{}, cmd string, args ...interface{}) Result {
	// db.SelectIn(&value, "select * from test  where a in (?)", [])
	// 传入切片的地址， 根据tag 的 db 自动补充，
	// 最求性能建议还是使用 GetRows or GetOne
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.Select(dest, newcmd, newargs...)
}

func makeArgs(cmd string, args ...interface{}) (string, []interface{}, error) {
	// 如果问号跟参数对不上， 报错
	count := strings.Count(cmd, "?")

	// vs: 新的参数

	if len(args) != count {
		return "", nil, fmt.Errorf("params error, expect %d, got %d", count, len(args))
	}
	vs := make([]interface{}, 0)
	// 找到？的索引
	need := make([]string, 0)
	cmdsplit := strings.Split(cmd, "?")
	for i, value := range args {
		typ := reflect.TypeOf(value)

		if typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice {
			invalue := reflect.ValueOf(args[i])
			if invalue.Len() == 0 {
				// 如果等于0，那么删除
				// 删除此cmd(的和下一个cmd的)
				end := strings.LastIndex(cmdsplit[i], "(")
				start := strings.Index(cmdsplit[i+1], ")")
				if end < 0 || start < 0 {
					return "", nil, errors.New("sql error")
				}
				cmdsplit[i+1] = cmdsplit[i+1][start+1:]
				// argsplit[value] = argsplit[value][:end]
				// 还要删除前面的in
				endin := strings.LastIndex(cmdsplit[i], "in")
				cmdsplit[i] = cmdsplit[i][:endin]
				fmt.Println(cmdsplit[i][:endin])
				cmdsplit[i] = strings.Trim(cmdsplit[i], " ")
				// 删除前面的word
				fmt.Println(cmdsplit[i])
				// 判断这个是不是not
				endspace := strings.LastIndex(cmdsplit[i], " ")
				if cmdsplit[i][endspace+1:] == "not" {
					cmdsplit[i] = cmdsplit[i][:endspace]
					cmdsplit[i] = strings.Trim(cmdsplit[i], " ")
					// 删除前面的word
					fmt.Println(cmdsplit[i])
					// 判断这个是不是not
					endspace = strings.LastIndex(cmdsplit[i], " ")
					cmdsplit[i] = cmdsplit[i][:endspace]
					cmdsplit[i] += " 1=1 "
				} else {
					cmdsplit[i] = cmdsplit[i][:endspace]
					cmdsplit[i] += " 1=0 "
				}

				// 删完后加上 1=0

				// 最后参数位置也要删掉
				// 判断是不是最后一个
				// 把当前的和下一个合并， 下一个的值为""
				need = append(need, "")
			} else {

				// 如果不为0的话，还是合并，不过添加了参数个的?
				wenhao := make([]string, invalue.Len())
				for i := 0; i < invalue.Len(); i++ {
					wenhao[i] = "?"
					vs = append(vs, invalue.Index(i).Interface())

				}
				need = append(need, strings.Join(wenhao, ","))
			}
			continue
		}
		need = append(need, "?")
		vs = append(vs, args[i])
		// 不是数组的话， 直接返回
	}
	// 去掉空值的
	cmds := ""
	for i := range need {
		cmds += cmdsplit[i] + " " + need[i] + " "
	}
	cmds += " " + cmdsplit[len(need)]
	return cmds, vs, nil
}
