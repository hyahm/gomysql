package main

import (
	"log"
)

type MeStruct struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type Person struct {
	ID        int64    `db:"id,default"`
	FirstName string   `db:"first_name,force"`
	LastName  string   `db:"last_name"`
	Email     string   `db:"email,force"`
	Me        MeStruct `db:"me"`
	Uids      []int64  `db:"uids"`
	TestJson  string   `db:"test"`
	Age       int      `json:"age" db:"age,counter"`
}

func main() {
	db, err := conf.NewMysqlDb()
	if err != nil {
		log.Fatal(err)
	}

	ps := &Person{
		FirstName: "what is it",
		LastName:  "hyahm.com",
		Email:     "aaaaa@eaml.com",
		// Me: MeStruct{
		// 	X: 10,
		// 	Y: 20,
		// 	Z: 30,
		// },
		// Uids: []int64{1},
		Age: 1,
	}
	// db.InsertInterfaceWithID(ps, "insert into person($key) values($value)")
	// $key  $value 是固定占位符
	// omitempty: 如果为空， 那么为数据库的默认值
	// struct, 指针， 切片 默认值为 ""
	// $set
	res := db.UpdateInterface(ps, "update person set $set where id=?", 1)
	if res.Err != nil {
		log.Fatal(res.Err)
	}
	db.Select()
	// cate := &User{}
	// res := db.Insert("INSERT INTO user (username, password) VALUES ('77tom', '123') ON DUPLICATE KEY UPDATE username='tom', password='123';")
	// // _, err = db.ReplaceInterface(&cate, "INSERT INTO user ($key) VALUES ($value) ON DUPLICATE KEY UPDATE $set")
	// if res.Err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(res.LastInsertId)
}
