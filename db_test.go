package gomysql

import "testing"

func TestSql(t *testing.T) {
	t.Log(ToPGSql("insert into account(username,password) values($1,$2)", "aaa", "bbb"))
}
