package main

import (
	"fmt"

	"github.com/hyahm/gomysql"
)

var (
	conf = &gomysql.Sqlconfig{
		Host:     "192.168.50.211",
		Port:     3306,
		UserName: "cander",
		Password: "123456",
		DbName:   "novel",
	}
)

func main() {

	db, err := conf.NewDb()
	if err != nil {
		panic(err)
	}
	var id int64
	err = db.GetOne("select id from shop_cover where id=? limit 1", 1).Scanf(&id)
	err = db.GetOne("select id from dp_book where id=? limit 1", 1).Scanf(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println(db.GetSql())
	fmt.Println(id)
}
