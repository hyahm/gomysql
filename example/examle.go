package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hyahm/gomysql"
)

var (
	conf = &gomysql.Sqlconfig{
		Host:         "192.168.50.250",
		Port:         3306,
		UserName:     "test",
		Password:     "123456",
		DbName:       "kaisa",
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
)

func main() {
	db, err := conf.NewDb()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 3)

	var id int64

	db.OpenDebug()
	err = db.GetOneIn("select id from subcategory where category_id in (?) and name in (?)", []string{"1"}, []string{"21", "23", "25"}).Scan(&id)
	fmt.Println(db.GetSql())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)

}
