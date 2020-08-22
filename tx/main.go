package main

import (
	"log"
	"os"
	"time"

	"github.com/hyahm/gomysql"
)

func main() {
	conf := &gomysql.Sqlconfig{
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
		LogFile:            ".failedlinux.sql",
	}
	db, err := conf.NewDb()
	if err != nil {
		os.Exit(1)
	}
	tx, err := db.NewTx(nil)
	_, err := tx.Insert("insert into test(name, age) value(?, ?)", "1", 2)
	if err != nil {
		log.Fatal(err)
	}
	twoid, err := tx.Insert("insert into test(name, age) value(?, ?)", "2", 4)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Update("update test set age=10 where id=?", twoid)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
}
