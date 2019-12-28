package main

import (
	"fmt"
	"github.com/hyahm/gomysql"
	"time"
)

var (
	conf = &gomysql.Sqlconfig{
		Host: "127.0.0.1",
		Port: 3306,
		UserName: "zth",
		Password: "123456789",
		DbName: "zth",
	}
)


func main() {
	var t time.Duration

	if t == time.Duration(0) {
		fmt.Println("1111")
	}
	//db, err := conf.NewDb()
	//if err != nil {
	//	panic(err)
	//}
	//var id int64
	//db.OpenDebug()
	//_, err = db.GetOne("select id from cmf_developer limit 1")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(db.GetSql())
	//fmt.Println(id)
}
