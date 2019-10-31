package gomysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
)

var (
	conf = &Sqlconfig{
		Host: "127.0.0.1",
		Port: 3306,
		UserName: "cander",
		Password: "123456",
		DbName: "aa",
	}
)


func Test_dbconn(t *testing.T) {

	// 保存配置并连接,  验证配置是否正确,  默认是连接db 对象的ping
	err := SaveConf("admin", conf)
	if err != nil {
		t.Fatal(err)
	}
	// 获取一个db 对象
	db, err := GetDb("admin")
	_,err = db.Insert("select * from test")
	if err != nil {
		t.Fatal(err)
	}
}

// 这是一个事物对象,  一样的使用
func Test_txconn(t *testing.T) {
	// 保存配置并连接,  验证配置是否正确
	err := SaveConf("admin", conf)
	if err != nil {
		t.Fatal(err)
	}
	// 获取一个db 对象
	tx, err := GetTx("admin")
	// 可以设置东西, 默认是这个
	// 测试断开后连接
	tx.Close()
	tx.Ctx = context.Background()
	rows,err := tx.GetRows("select * from test")
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
	// 保存配置并连接,  验证配置是否正确
	err := SaveConf("admin", conf)
	if err != nil {
		t.Fatal(err)
	}
	// 获取一个db 对象
	tx, err := GetTx("admin")
	// 可以设置东西, 默认是这个
	// 测试断开后连接
	tx.Close()

	id,err := tx.Insert("insert into test(id, name) values(?,?)", 675544677, "password")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)


	err = tx.Commit()
	if err != nil {
		t.Log(err)
	}
	//var user string
	//err = tx.GetOne("select name from test where id=?", id).Scan(&user)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Logf("username: %s", user)
}


// 测试提交
func TestTx(t *testing.T) {
	// 保存配置并连接,  验证配置是否正确
	err := SaveConf("admin", conf)
	if err != nil {
		t.Fatal(err)
	}
	// 获取一个db 对象, 已经begin
	tx, err := GetTx("admin")
	// 可以设置东西, 默认是这个
	// 测试断开后连接
	tx.Begin()
	id,err := tx.Insert("insert into test(id, name) values(?,?)", 1954487, "password")
	if err != nil {
		t.Log(err)
	}
	t.Log(id)
	tx.Commit()
	var user string
	tx.Begin()
	err = tx.GetOne("select name from test where id=?", 1954487).Scan(&user)
	if err != nil {
		t.Fatal(err) // 这里会出现没有找到行
	}
	t.Logf("username: %s", user)
}

func BenchmarkInsert(b *testing.B) {
	err := SaveConf("admin", conf)
	if err != nil {
		b.Fatal(err)
	}
	// 获取一个db 对象, 已经begin
	d, err := GetDb("admin")
	for i:= 0; i < b.N; i++ {
		_, err = d.Insert("insert into test(id,name) values(?,?)",i, fmt.Sprintf("name %d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkLocalInsert(b *testing.B) {
	//err := SaveConf("admin", conf)
	//if err != nil {
	//	b.Fatal(err)
	//}
	//// 获取一个db 对象, 已经begin
	connstring := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		conf.UserName, conf.Password, conf.Host, conf.Port, conf.DbName,
	)
	conn, err := sql.Open("mysql", connstring)
	if err != nil {
		log.Fatal(err)
	}

	//d, err := GetDb("admin")
	for i:= 0; i < b.N; i++ {
		_, err = conn.Exec("insert into test(id,name) values(?,?)",i, fmt.Sprintf("name %d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
}



