package main

import (
	"fmt"
	"github.com/hyahm/gomysql"
)

var (
	conf = &gomysql.Sqlconfig{
		Host: "127.0.0.1",
		Port: 3306,
		UserName: "zth",
		Password: "123456",
		DbName: "zth",
	}
)


func main() {
	db, err := conf.NewDb()
	if err != nil {
		panic(err)
	}
	var id int64
	db.OpenDebug()
	err = db.GetOne("select id from cmf_developer limit 1").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println(db.PrintSql())
	fmt.Println(id)
}
