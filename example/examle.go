package main

import (
	"fmt"
	"log"

	"github.com/hyahm/gomysql"
)

var (
	conf = &gomysql.Sqlconfig{
		Host:         "192.168.0.107",
		Port:         3306,
		UserName:     "cander",
		Password:     "123456",
		DbName:       "bug",
		MaxOpenConns: 1,
	}
)

func main() {

	db, err := conf.NewDb()
	if err != nil {
		log.Fatal(err)
	}
	var id int64
	db.OpenDebug()
	db.GetOneIn("select * from xxx where id=? and a in (?) and b in (?) and name=?",
		6666,
		(gomysql.InArgs)([]string{"1", "2", "4", "5", "6", "7", "8", "89", "3", "4"}).ToInArgs(),
		(gomysql.InArgs)([]string{"aaa", "bbb"}).ToInArgs(),
		"cander").Scan(&id)
	fmt.Println(db.GetSql())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(id)
	// _, err = db.GetRows("select id from dp_book")
	// if err != nil {
	// 	panic(err)
	// }
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
