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
		DbName:       "xilin",
		MaxOpenConns: 10,
		MaxIdleConns: 1,
	}
)

func main() {
	db, err := conf.NewDb()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 3)

	var id int64

	countArgs := make([]interface{}, 0)
	statuslist := []string{"12"}
	myproject := []string{"21", "23", "25"}
	countArgs = append(countArgs, (gomysql.InArgs)(statuslist).ToInArgs())
	countArgs = append(countArgs, (gomysql.InArgs)(myproject).ToInArgs())
	db.OpenDebug()
	err = db.GetOneIn("select id from user a in (?) and b in (?)", countArgs...).Scan(&id)
	fmt.Println(db.GetSql())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)

}
