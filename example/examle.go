package main

import (
	"fmt"
	"log"

	"github.com/hyahm/gomysql"
)

var (
	conf = &gomysql.Sqlconfig{
		Host:         "192.168.50.211",
		Port:         3306,
		UserName:     "cander",
		Password:     "123456",
		DbName:       "novel",
		MaxOpenConns: 1,
	}
)

func main() {

	db, err := conf.NewDb()
	if err != nil {
		panic(err)
	}
	var id int64
	err = db.GetOne("select id from dp_book").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
	_, err = db.GetRows("select id from dp_book")
	if err != nil {
		panic(err)
	}
	// rows.Close()

	// rows.Close()
	// var id int
	// for rows.Next() {

	// 	err = rows.Scan(&id)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	fmt.Println(id)
	// }
}
