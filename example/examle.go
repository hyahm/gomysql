package main

import (
	"fmt"
	"time"

	"github.com/hyahm/gomysql"
)

var (
	conf = &gomysql.Sqlconfig{
		Host:     "172.247.132.58",
		Port:     3306,
		UserName: "root",
		Password: "$h@F6ph6&UdY",
		DbName:   "shop",
	}
)

func main() {
	var t time.Duration

	if t == time.Duration(0) {
		fmt.Println("1111")
	}
	db, err := conf.NewDb()
	if err != nil {
		panic(err)
	}
	var id int64
	db.OpenDebug()
	_, err = db.GetOne("select id from shop_cover where id=? limit 1", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(db.GetSql())
	fmt.Println(id)
}
