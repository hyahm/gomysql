package main

import (
	"fmt"
	"log"

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

	db, err := conf.NewDb()
	if err != nil {
		log.Fatal(err)
	}

	db.NewCurder(User{})
	// cate := &User{}
	res := db.Insert("INSERT INTO user (username, password) VALUES ('77tom', '123') ON DUPLICATE KEY UPDATE username='tom', password='123';")
	// _, err = db.ReplaceInterface(&cate, "INSERT INTO user ($key) VALUES ($value) ON DUPLICATE KEY UPDATE $set")
	if res.Err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.LastInsertId)
}
