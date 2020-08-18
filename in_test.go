package gomysql

import (
	"testing"
)

func TestInMany(t *testing.T) {
	cmd, args, err := makeArgs("select * from xxx where id=? and a in (?)", 6666, []interface{}{1, 2, 4, 5, 6, 7, 8, 89, 3, 4})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cmd)
	t.Log(args)
}

func TestInOne(t *testing.T) {
	cmd, args, err := makeArgs("select * from xxx where id=? and a in (?)", 6666, []interface{}{1})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cmd)
	t.Log(args)
}

func TestInEmpty(t *testing.T) {
	cmd, args, err := makeArgs("select * from xxx where id=? and a in (?)", 6666, []interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cmd)
	t.Log(args)
}

// func TestReplace(t *testing.T) {
// 	// 讲索引为什么的？。 替换成n个
// 	am, err := replace("select * from xxx where id=? and a in (?)", 1, 10)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(am)
// }
