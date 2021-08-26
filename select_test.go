package gomysql

import (
	"encoding/json"
	"testing"
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
	db, err := conf.NewDb()
	if err != nil {
		t.Fatal(err)
	}
	// db.Query(schema)
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "John", "Doe", "johndoeDNE@gmail.net")
	persons := &Person{}
	err = db.Select(&persons, "select * from person limit 1")
	if err != nil {
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
	db, err := conf.NewDb()
	if err != nil {
		t.Fatal(err)
	}
	// db.Query(schema)
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "John", "Doe", "johndoeDNE@gmail.net")
	persons := make([]Person, 0)
	err = db.Select(&persons, "select * from person limit 1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(persons)
	// 建表

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
	db, err := conf.NewDb()
	if err != nil {
		t.Fatal(err)
	}
	// db.Query(schema)
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "John", "Doe", "johndoeDNE@gmail.net")
	persons := make([]*Person, 0)
	err = db.Select(&persons, "select * from person limit 1")
	if err != nil {
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
	db, err := conf.NewDb()
	if err != nil {
		t.Fatal(err)
	}
	// db.Query(schema)
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// db.Insert("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", "John", "Doe", "johndoeDNE@gmail.net")
	persons := Person{}
	err = db.Select(&persons, "select * from person limit 1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(persons)
	// 建表

}
