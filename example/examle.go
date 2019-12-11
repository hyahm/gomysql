package main

import (
	"fmt"
	"github.com/hyahm/gomysql"
)

var (
	conf = &gomysql.Sqlconfig{
		Host: "120.24.171.222",
		Port: 3306,
		UserName: "zth",
		Password: "f^9^NgW2WszEUB%P",
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
	fmt.Printf("%+v \n", db)
	err = db.GetOne("select id from cmf_developer limit 1").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println(db.PrintSql())
	fmt.Println(id)
}
