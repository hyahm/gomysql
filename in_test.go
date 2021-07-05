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

func TestString(t *testing.T) {
	a := "select * from xxx where id=? and a in (?) and b in (?) and name=?"
	args := []interface{}{
		6666,
		[]string{"1", "2", "4", "5", "6", "7", "8", "89", "3", "4"},
		[]string{},
		"cander",
	}
	cmd, args, err := makeArgs(a, args...)
	if err != nil {
		t.Log(err)
	}
	t.Log(cmd)
	t.Log(args)
}

func TestInArgs(t *testing.T) {
	td := []testData{
		{
			title: "参数都有的 测试前后都有条件的",
			cmd:   "select * from xxx where id=? and a in (?) and b in (?) and name=?",
			args: []interface{}{
				6666,
				[]string{"1", "2", "4", "5", "6", "7", "8", "89", "3", "4"},
				[]string{},
				"cander",
			},
			expectCmd:  "select * from xxx where id=? and a in (?,?,?,?,?,?,?,?,?,?) and name=?",
			expectArgs: []interface{}{6666, "1", "2", "4", "5", "6", "7", "8", "89", "3", "4", "cander"},
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
	t.Log(cmd)
	t.Log(args)
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
