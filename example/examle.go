package main

import (
	"log"

	"github.com/hyahm/golog"
	"github.com/hyahm/gomysql"
)

type User struct {
	Id       int64  `json:"id" db:"id,omitempty"`
	Username string `json:"username" db:"username"` // 分类英文名， 文件夹命名 唯一索引
	Password string `json:"password" db:"password"`
}

type Postparam struct {
}

var (
	conf = &gomysql.Sqlconfig{
		Host:         "192.168.50.250",
		Port:         3306,
		UserName:     "test",
		Password:     "123456",
		DbName:       "test",
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
)

type MeStruct struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type Person struct {
	ID        int64    `db:"id,default"`
	FirstName string   `db:"first_name,force"`
	LastName  string   `db:"last_name"`
	Email     string   `db:"email,force"`
	Me        MeStruct `db:"me"`
	Uids      []int64  `db:"uids"`
	TestJson  string   `db:"test"`
	Age       int      `json:"age" db:"age,counter"`
}

func main() {
	defer golog.Sync()
	db, err := conf.NewDb()
	if err != nil {
		log.Fatal(err)
	}

	ps := &Person{
		FirstName: "what is it",
		LastName:  "hyahm.com",
		Email:     "aaaaa@eaml.com",
		// Me: MeStruct{
		// 	X: 10,
		// 	Y: 20,
		// 	Z: 30,
		// },
		Uids: []int64{1},
		Age:  1,
	}

	// $key  $value 是固定占位符
	// omitempty: 如果为空， 那么为数据库的默认值
	// struct, 指针， 切片 默认值为 ""
	// $set
	res := db.UpdateInterface(ps, "update person set $set where id=?", 4)
	golog.Info(res.Sql)
	if res.Err != nil {
		golog.Fatal(res.Err)
	}

	// cate := &User{}
	// res := db.Insert("INSERT INTO user (username, password) VALUES ('77tom', '123') ON DUPLICATE KEY UPDATE username='tom', password='123';")
	// // _, err = db.ReplaceInterface(&cate, "INSERT INTO user ($key) VALUES ($value) ON DUPLICATE KEY UPDATE $set")
	// if res.Err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(res.LastInsertId)
}
