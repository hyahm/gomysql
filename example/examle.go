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
		DbName:       "test",
		MaxOpenConns: 1,
	}
)

func main() {

	db, err := conf.NewDb()
	if err != nil {
		panic(err)
	}
	s, err := db.InsertMany("insert into test(name, age) values(?,?)", "test3", 11, "test4", 2, "test3", 5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
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
