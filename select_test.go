package gomysql

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

type Paipan struct {
	ID          int64           `db:"id"`
	Name        string          `db:"name"`
	Gender      bool            `db:"gender"`
	Mark        string          `db:"mark"`
	Data        json.RawMessage `db:"data"`
	LifeStyleId int             `db:"life_style_id"`
	Created     int64           `db:"created"`
}

type Category struct {
	ID      int64   `db:"id"`
	Uids    []int64 `db:"uids"`
	Subcate []int64 `db:"subcate"`
}

var schema = `
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);

CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
)`

func TestSelect(t *testing.T) {
	conf := Sqlconfig{
		UserName:        "test",
		Password:        "123456",
		Port:            3306,
		DbName:          "test",
		Host:            "192.168.101.4",
		MultiStatements: true,
	}
	db, err := conf.NewMysqlDb()
	if err != nil {
		t.Fatal(err)
	}
	// db.Query(schema)
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "John", "Doe", "johndoeDNE@gmail.net")
	persons := &Person{}
	res := db.Select(&persons, "select * from person limit 1")
	if res.Err != nil {
		t.Fatal(err)
	}
	t.Log(persons)
	// 建表

}

func TestSelect1(t *testing.T) {
	conf := Sqlconfig{
		UserName:        "test",
		Password:        "123456",
		Port:            3306,
		DbName:          "test",
		Host:            "192.168.101.4",
		MultiStatements: true,
	}
	db, err := conf.NewMysqlDb()
	if err != nil {
		t.Fatal(err)
	}
	// db.Query(schema)
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "John", "Doe", "johndoeDNE@gmail.net")
	persons := make([]Person, 0)
	res := db.Select(&persons, "select * from person limit 1")
	if res.Err != nil {
		t.Fatal(err)
	}
	t.Log(persons)
	// 建表

}

const Num = 1000

func TestMysql8(t *testing.T) {
	start := time.Now()
	conf := &Sqlconfig{
		Host:               "127.0.0.1",
		UserName:           "cander",
		Password:           "123456",
		DbName:             "test",
		Port:               3306,
		MaxOpenConns:       100,
		MaxIdleConns:       10,
		ReadTimeout:        100 * time.Second,
		WriteTimeout:       100 * time.Second,
		WriteLogWhenFailed: true,
		ConnMaxLifetime:    30 * time.Second,
		// 删改查失败写入的文件
		LogFile: ".failedlinux.sql",
	}
	ch := make(chan int, Num)
	db, err := conf.NewMysqlDb()
	if err != nil {
		os.Exit(1)
	}

	for i := 0; i < Num; i++ {
		go func(i int) {
			db.Insert("insert into test(name, age) values(?,?)", fmt.Sprintf("test%d", i), i)
			ch <- 1
		}(i)

	}

	for i := 0; i < Num; i++ {
		<-ch
	}

	rows, err := db.GetRowsIn("select id from test where age in (?)", []string{"1", "2", "3", "4", "5"})
	if err != nil {
		os.Exit(1)
	}
	for rows.Next() {
		var id int64
		rows.Scan(&id)
	}
	rows.Close()
	log.Println("mysql8:", time.Since(start).Seconds())
}

func TestSelect2(t *testing.T) {
	conf := Sqlconfig{
		UserName:        "test",
		Password:        "123456",
		Port:            3306,
		DbName:          "test",
		Host:            "192.168.101.4",
		MultiStatements: true,
	}
	db, err := conf.NewMysqlDb()
	if err != nil {
		t.Fatal(err)
	}
	// db.Query(schema)
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "John", "Doe", "johndoeDNE@gmail.net")
	persons := make([]*Person, 0)
	res := db.Select(&persons, "select * from person limit 1")
	if res.Err != nil {
		t.Fatal(err)
	}
	t.Log(*persons[0])
	// 建表

}

func TestSelect3(t *testing.T) {
	conf := Sqlconfig{
		UserName:        "test",
		Password:        "123456",
		Port:            3306,
		DbName:          "test",
		Host:            "192.168.101.4",
		MultiStatements: true,
	}
	db, err := conf.NewMysqlDb()
	if err != nil {
		t.Fatal(err)
	}
	// db.Query(schema)
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "John", "Doe", "johndoeDNE@gmail.net")
	persons := Person{}
	res := db.Select(&persons, "select * from person limit 1")
	if res.Err != nil {
		t.Fatal(err)
	}
	t.Log(persons)
	// 建表

}
