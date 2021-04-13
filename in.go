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

func makeArgs(cmd string, args ...interface{}) (string, []interface{}, error) {
	// 如果问号跟参数对不上， 报错
	count := strings.Count(cmd, "?")

	vs := make([]interface{}, 0)
	if len(args) != count {
		return "", vs, fmt.Errorf("params error, expect %d, got %d", count, len(args))
	}
	var err error
	inIndex := 0
	for _, value := range args {
		typ := reflect.TypeOf(value)
		vv := reflect.ValueOf(value)

		if typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice {
			l := vv.Len()
			if l == 0 {
				vs = append(vs, "")
			} else {
				cmd, err = replace(cmd, inIndex, l)
				if err != nil {
					return cmd, vs, err
				}
				for i := 0; i < l; i++ {
					vs = append(vs, vv.Index(i).Interface())
				}

				inIndex++
			}

		} else {
			vs = append(vs, value)
		}

		// 不是数组的话， 直接返回
	}
	return cmd, vs, nil
}

// func findStrIndex(cmd string, pos int, del bool) (string, error) {
// 	count := strings.Count(cmd, "?")
// 	tmp := cmd
// 	lastcmd := ""
// 	preStr := ""
// 	sufStr := ""
// 	start := 0
// 	for i := 0; i < count; i++ {
// 		thisIndex := strings.Index(tmp, "?")
// 		start = start + thisIndex + 1
// 		tmp = tmp[thisIndex+1:]

// 		if i == pos {
// 			// 找到前面的(
// 			ksindex := strings.Index(cmd[:start], "(")
// 			klindex := strings.Index(cmd[start:], ")")

// 			if ksindex < 0 || klindex < 0 {
// 				return "", errors.New("sql error")
// 			}
// 			preStr = cmd[:ksindex]
// 			sufStr = cmd[klindex+start+1:]
// 			inIndex := strings.LastIndex(preStr, "in")
// 			if inIndex < 0 {
// 				return "", errors.New("not found in , please do not use in Func")
// 			}

// 			// 删除in
// 			preStr = strings.Trim(cmd[:inIndex], " ")
// 			if del {
// 				// 找到前面一个空格位置
// 				// 去掉前面的单词
// 				beforeIn, word := getLastStr(preStr)
// 				if word == "not" {
// 					// 如果前面是not， 那么还要去掉一次
// 					beforeIn, _ = getLastStr(beforeIn)
// 				}
// 				preStr = beforeIn
// 				//  end 删除前面的
// 				sufStr = strings.Trim(cmd[klindex+start+1:], " ")
// 				deleteCondition := true
// 				if sufStr != "" {
// 					// 不为空先删除or 或者 and
// 					conditionStr, word := getNextStr(sufStr)
// 					if word == "or" || word == "and" {
// 						// 删除后面的
// 						sufStr = conditionStr
// 						deleteCondition = false
// 					}
// 				}
// 				if deleteCondition {
// 					beforeIn, word := getLastStr(preStr)
// 					if word == "where" || word == "on" {
// 						preStr = beforeIn
// 					}
// 				}
// 				// 删除前面的 tiaonian
// 				if sufStr == "" {
// 					// 如果后面没东西， 还要删除 and or
// 					beforeIn, word := getLastStr(preStr)
// 					if word == "and" || word == "or" {
// 						preStr = beforeIn
// 					}

// 				}

// 				lastcmd = preStr + sufStr

// 			} else {
// 				beforeIn, word := getLastStr(preStr)
// 				if word == "not" {
// 					// 如果前面是not， 那么还要去掉一次
// 					preStr = beforeIn
// 					lastcmd = preStr + "<>?" + sufStr
// 				} else {
// 					lastcmd = preStr + "=?" + sufStr
// 				}
// 			}

// 			break

// 		}
// 	}

// 	return lastcmd, nil
// }

func replace(cmd string, index int, count int) (string, error) {
	// 替换？,
	// index: 替换第几个in的？
	// count: 替换多少次

	// 先寻找in的位置, 然后寻找后面的?
	tmp := strings.ToLower(cmd)
	m := make([]string, count)
	for j := 0; j < count; j++ {
		m[j] = "?"
	}

	c := strings.Count(tmp, " in ")
	if index > c-1 {
		return "", ErrRangeOut
	}
	start := 0
	for i := 0; i < c; i++ {
		thisInIndex := strings.Index(tmp[start:], " in ")
		// 找到后面第一个?
		start += thisInIndex + 4
		if index == i {
			thisIndex := strings.Index(tmp[start:], "?")
			start += thisIndex + 1
			return tmp[:start-1] + strings.Join(m, ",") + tmp[start:], nil
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

func (d *Db) GetOneIn(cmd string, args ...interface{}) *Row {
	newcmd, newargs, err := makeArgs(cmd, args...)
	if err != nil {
		panic(err)
	}
	return d.GetOne(newcmd, newargs...)
}

// func getLastStr(s string) (string, string) {
// 	// 通过一个字符串， 获取新的字符串 上一个空格字符的上一个word， 以及字符串中的位置
// 	new := strings.Trim(s, " ")
// 	index := strings.LastIndex(new, " ")
// 	word := new[index+1:]
// 	return new[:index], word
// }

// func getNextStr(s string) (string, string) {
// 	// 通过一个字符串， 截取后的新的字符串 上一个空格字符的上一个word， 以及字符串中的位置
// 	new := strings.Trim(s, " ")
// 	index := strings.Index(new, " ")
// 	word := new[:index]
// 	return new[index:], word
// }
