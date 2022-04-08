package main

import (
	"fmt"
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

func main() {
	defer golog.Sync()
	db, err := conf.NewDb()
	if err != nil {
		log.Fatal(err)
	}

	curd := db.NewCurder("user")
	user := make([]User, 0)
	user = append(user, User{
		Username: "jack1",
		Password: "97979",
	})

	user = append(user, User{
		Username: "jack2",
		Password: "979792",
	})

	res := curd.Create(user)

	if res.Err != nil {
		log.Fatal(res.Err)
	}
	fmt.Println(res.Sql)
	fmt.Println(res.LastInsertId)
	fmt.Println(res.LastInsertIds)
	// cate := &User{}
	// res := db.Insert("INSERT INTO user (username, password) VALUES ('77tom', '123') ON DUPLICATE KEY UPDATE username='tom', password='123';")
	// // _, err = db.ReplaceInterface(&cate, "INSERT INTO user ($key) VALUES ($value) ON DUPLICATE KEY UPDATE $set")
	// if res.Err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(res.LastInsertId)
}
