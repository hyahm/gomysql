package main

import (
	"fmt"
	"gomysql"
	"log"
)

func main() {
	Conf := &gomysql.Sqlconfig{
		UserName: "root",
		Password: "123456",
		Port: 3306,
		DbName: "admin",
	}
	gomysql.SaveConf("x7", Conf)
	gomysql.ConnDB("x7")
	rows,err := gomysql.GetRows("x7", "select username,password from user")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next(){
		fmt.Println("11111111111111111")
		var user, pwd string
		err = rows.Scan(&user, &pwd)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("username: %s, password: %s \n", user, pwd)
	}
}
