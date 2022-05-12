package main

import (
	"fmt"
	"os"

	"github.com/hyahm/gomysql"
)

type Account struct {
	Id       int    `json:"id" db:"id,omitempty"`
	Username string `json:"username" db:"username"` // 分类英文名， 文件夹命名 唯一索引
	Password string `json:"password" db:"password"`
	Gender   bool   `json:"gender" db:"gender"`
}

var (
	conf = &gomysql.Sqlconfig{
		Host:     "192.168.101.13",
		Port:     5432,
		UserName: "test",
		Password: "123456",
		DbName:   "mydb",
		Debug:    true,
	}
)

func main() {
	pg, err := conf.NewPGPool()
	// urlExample := "postgres://test:123456@192.168.50.58:5432/mydb"
	// conn, err := pgxpool.Connect(context.Background(), "postgres://test:123456@192.168.50.58:5432/mydb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// var id int64
	// row := pg.QueryRow(context.Background(), "insert into account(username, password) values($1, $2) returning id", "Aaa", "bbb")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// row.Scan(&id)
	// fmt.Println(id)
	addAccount := Account{
		Username: "99999",
		Gender:   true,
	}

	res := pg.InsertInterfaceWithID(addAccount, "insert into account($key) values($value) returning id")
	fmt.Println(res.Sql)
	fmt.Println(res)

	// res := pg.UpdateInterface(account, "update account set $set where id=11")
	// fmt.Println(res.Sql)
	// fmt.Println(res)
	// fmt.Println(greeting)
	// var name string
	// var weight int64
	// err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }

	account := make([]Account, 0)
	pg.Select(&account, "select * from account")
	fmt.Println(account)
	// conn.QueryRow()
	// fmt.Println(name, weight)
}
