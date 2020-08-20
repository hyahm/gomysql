package gomysql

import (
	"reflect"
	"strings"
	"testing"
)

type testData struct {
	title      string
	cmd        string
	args       []interface{}
	expectCmd  string
	expectArgs []interface{}
}

func TestIn(t *testing.T) {
	td := []testData{
		{
			title: "参数都有的 测试前后都有条件的",
			cmd:   "select * from xxx where id=? and a in (?) and name=?",
			args: []interface{}{
				6666,
				(InArgs)([]string{"1", "2", "4", "5", "6", "7", "8", "89", "3", "4"}).ToInArgs(),
				"cander",
			},
			expectCmd:  "select * from xxx where id=? and a in (?,?,?,?,?,?,?,?,?,?) and name=?",
			expectArgs: []interface{}{6666, 1, 2, 4, 5, 6, 7, 8, 89, 3, 4, "cander"},
		},
		{
			title: "参数只有一个， 前面有条件，后面也有条件",
			cmd:   "select * from xxx where id=? and a in (?) and name=?",
			args: []interface{}{
				6666,
				[]interface{}{1},
				"cander",
			},
			expectCmd:  "select * from xxx where id=? and a=? and name=?",
			expectArgs: []interface{}{6666, 1, "cander"},
		},
		{
			title: "参数只有一个， not ",
			cmd:   "select * from xxx where id=? and a not in (?) and name=?",
			args: []interface{}{
				6666,
				[]interface{}{1},
				"cander",
			},
			expectCmd:  "select * from xxx where id=? and a<>? and name=?",
			expectArgs: []interface{}{6666, 1, "cander"},
		},
		{
			title: "参数空的，测试前面有条件，后面没有条件",
			cmd:   "select * from xxx where id=? and a in (?)",
			args: []interface{}{
				6666,
				[]interface{}{},
			},
			expectCmd:  "select * from xxx where id=?",
			expectArgs: []interface{}{6666},
		},
		{
			title: "参数空的，测试前面有条件，后面有条件",
			cmd:   "select * from xxx where a in (?) and name=?",
			args: []interface{}{
				[]interface{}{},
				"cander",
			},
			expectCmd:  "select * from xxx where name=?",
			expectArgs: []interface{}{"cander"},
		},
		{
			title: "参数空的，测试前面没有条件，后面没有条件",
			cmd:   "select * from xxx where a in (?)",
			args: []interface{}{
				[]interface{}{},
			},
			expectCmd:  "select * from xxx",
			expectArgs: []interface{}{},
		},
		{
			title: "参数空的，测试前面没有条件，后面有条件",
			cmd:   "select * from xxx where a in (?) and name=?",
			args: []interface{}{
				[]interface{}{},
				"cander",
			},
			expectCmd:  "select * from xxx where name=?",
			expectArgs: []interface{}{"cander"},
		},
	}
	for _, v := range td {
		run(t, v)
	}
}

func run(t *testing.T, td testData) {

	cmd, args, err := makeArgs(td.cmd, td.args...)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Trim(cmd, " ") == td.expectCmd {
	} else {
		t.Log("title", td.title)
		t.Log("cmd failed")
	}
	if reflect.DeepEqual(args, td.expectArgs) {
	} else {
		t.Log("title", td.title)
		t.Log("args failed")
	}

}

// func TestReplace(t *testing.T) {
// 	// 讲索引为什么的？。 替换成n个
// 	am, err := replace("select * from xxx where id=? and a in (?)", 1, 10)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(am)
// }
