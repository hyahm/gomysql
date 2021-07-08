package gomysql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// RangeOutErr 错误信息
var ErrRangeOut = errors.New("索引值超出范围")

// IgnoreWords 如果删除in的字段， 前面如果出现下面的关键字， 也删除
var IgnoreWords = []string{"or", "and"}

// 如果后面有or 或者 and，只用删除or 或者and 那么删除前面的where， on 关键字
var IgnoreCondition = []string{"where", "on"}

// 如果后面有or 或者 and，只用删除or 或者and
// 如果没有， 那么删除前面的 where 或者 on

// func makeArgs(cmd string, args ...interface{}) (string, []interface{}, error) {
// 	// 如果问号跟参数对不上， 报错
// 	count := strings.Count(cmd, "?")

// 	// vs: 新的参数
// 	vs := make([]interface{}, 0)
// 	if len(args) != count {
// 		return "", vs, fmt.Errorf("params error, expect %d, got %d", count, len(args))
// 	}
// 	var err error
// 	inIndex := 0
// 	// 找到in的索引
// 	for _, value := range args {
// 		typ := reflect.TypeOf(value)
// 		vv := reflect.ValueOf(value)

// 		if typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice {
// 			l := vv.Len()
// 			if l == 0 {
// 				// 删除2边括号和前面2个空格的字符
// 				// 先删除2边的括号
// 				// 左边：

// 				fmt.Println(vv.Interface())
// 				vs = append(vs, "")
// 			} else {
// 				cmd, err = replace(cmd, inIndex, l)
// 				if err != nil {
// 					return cmd, vs, err
// 				}
// 				for i := 0; i < l; i++ {
// 					vs = append(vs, vv.Index(i).Interface())
// 				}

// 				inIndex++
// 			}

// 		} else {
// 			vs = append(vs, value)
// 		}

// 		// 不是数组的话， 直接返回
// 	}
// 	fmt.Println(cmd)
// 	fmt.Println(vs)
// 	return cmd, vs, nil
// }

// func replace(cmd string, index int, count int) (string, error) {
// 	// 替换？,
// 	// index: 替换第几个in的？
// 	// count: 替换多少次

// 	// 先寻找in的位置, 然后寻找后面的?
// 	tmp := strings.ToLower(cmd)
// 	m := make([]string, count)
// 	for j := 0; j < count; j++ {
// 		m[j] = "?"
// 	}

// 	c := strings.Count(tmp, " in ")
// 	if index > c-1 {
// 		return "", ErrRangeOut
// 	}
// 	start := 0
// 	for i := 0; i < c; i++ {
// 		thisInIndex := strings.Index(tmp[start:], " in ")
// 		// 找到后面第一个?
// 		start += thisInIndex + 4
// 		if index == i {
// 			thisIndex := strings.Index(tmp[start:], "?")
// 			start += thisIndex + 1
// 			return tmp[:start-1] + strings.Join(m, ",") + tmp[start:], nil
// 		}
// 	}
// 	return cmd, nil
// }

func (d *Db) UpdateIn(cmd string, args ...interface{}) (int64, error) {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return 0, err
	}
	return d.Update(newcmd, newargs...)
}

func (d *Db) InsertIn(cmd string, args ...interface{}) (int64, error) {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return 0, err
	}
	return d.Insert(newcmd, newargs...)
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

func (d *Db) UpdateInterfaceIn(dest interface{}, cmd string, args ...interface{}) (int64, error) {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return 0, err
	}
	return d.UpdateInterface(dest, newcmd, newargs...)
}

func (d *Db) InsertInterfaceWithoutIDIn(dest interface{}, cmd string, args ...interface{}) error {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return err
	}
	return d.InsertInterfaceWithoutID(dest, newcmd, newargs...)
}

func (d *Db) InsertInterfaceWithIDIn(dest interface{}, cmd string, args ...interface{}) ([]int64, error) {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return nil, err
	}
	return d.InsertInterfaceWithID(dest, newcmd, newargs...)
}

func (d *Db) SelectIn(dest interface{}, cmd string, args ...interface{}) error {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		return err
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

				cmdsplit[i] = strings.Trim(cmdsplit[i], " ")
				// 删除前面的word
				endspace := strings.LastIndex(cmdsplit[i], " ")
				cmdsplit[i] = cmdsplit[i][:endspace]
				// 删完后加上 1=0
				cmdsplit[i] += " 1=0 "

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
