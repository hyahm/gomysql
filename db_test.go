package gomysql

import (
	"context"
	"testing"
)

func Test_dbconn(t *testing.T) {
	c := &Sqlconfig{
		Host: "127.0.0.1",
		Port: 3306,
		UserName: "root",
		Password: "123456",
		DbName: "admin",
	}
	// 保存配置并连接,  验证配置是否正确,  默认是连接db 对象的ping
	err := SaveConf("admin", c)
	if err != nil {
		t.Fatal(err)
	}
	// 获取一个db 对象
	db, err := GetDb("admin")
	_,err = db.Insert("select * from user")
	if err != nil {
		t.Fatal(err)
	}
}

// 这是一个事物对象,  一样的使用
func Test_txconn(t *testing.T) {
	c := &Sqlconfig{
		Host: "127.0.0.1",
		Port: 3306,
		UserName: "root",
		Password: "123456",
		DbName: "x7",
	}
	// 保存配置并连接,  验证配置是否正确
	err := SaveConf("admin", c)
	if err != nil {
		t.Fatal(err)
	}
	// 获取一个db 对象
	tx, err := GetTx("admin")
	// 可以设置东西, 默认是这个
	// 测试断开后连接
	tx.Close()
	tx.Ctx = context.Background()
	rows,err := tx.GetRows("select * from user")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("22222")
	for rows.Next() {
		t.Log("111")
		t.Log(rows.Columns())
	}

}

// 测试提交
func Test_txcommit(t *testing.T) {
	c := &Sqlconfig{
		Host: "127.0.0.1",
		Port: 3306,
		UserName: "root",
		Password: "123456",
		DbName: "x7",
	}
	// 保存配置并连接,  验证配置是否正确
	err := SaveConf("admin", c)
	if err != nil {
		t.Fatal(err)
	}
	// 获取一个db 对象
	tx, err := GetTx("admin")
	// 可以设置东西, 默认是这个
	// 测试断开后连接
	tx.Close()

	id,err := tx.Insert("insert into user(username, password) values(?,?)", "username", "password")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)


	err = tx.Commit()
	if err != nil {
		t.Log(err)
	}
	var user string
	err = tx.GetOne("select username from user where id=?", id).Scan(&user)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("username: %s", user)
}


// 测试提交
func Test_txnotcommit(t *testing.T) {
	c := &Sqlconfig{
		Host: "127.0.0.1",
		Port: 3306,
		UserName: "root",
		Password: "123456",
		DbName: "x7",
	}
	// 保存配置并连接,  验证配置是否正确
	err := SaveConf("admin", c)
	if err != nil {
		t.Fatal(err)
	}
	// 获取一个db 对象, 已经begin
	tx, err := GetTx("admin")
	// 可以设置东西, 默认是这个
	// 测试断开后连接
	tx.Close()

	id,err := tx.Insert("insert into user(username, password) values(?,?)", "username", "password")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

	var user string
	err = tx.GetOne("select username from user where id=?", id).Scan(&user)
	if err != nil {
		t.Fatal(err) // 这里会出现没有找到行
	}
	t.Logf("username: %s", user)
}

