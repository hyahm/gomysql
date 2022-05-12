package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hyahm/gomysql"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go Insert8(wg)
	// go Insert5(wg)
	wg.Wait()
}

// Num 插入的次数
const Num = 1000

// Insert8 mysql8 的插入
func Insert8(wg *sync.WaitGroup) {

	start := time.Now()
	conf := &gomysql.Sqlconfig{
		Host:               "192.168.50.211",
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
	log.Println("mysql8:", time.Since(start).Seconds())
	wg.Done()
}

// Insert5 mysql5 的插入
func Insert5(wg *sync.WaitGroup) {
	start := time.Now()
	conf := &gomysql.Sqlconfig{
		Host:               "127.0.0.1",
		UserName:           "root",
		Password:           "123456",
		DbName:             "test",
		Port:               3306,
		MaxOpenConns:       50000,
		MaxIdleConns:       50000,
		ConnMaxLifetime:    time.Minute * 10,
		ReadTimeout:        time.Second * 10,
		WriteLogWhenFailed: true,
		LogFile:            ".failedwindows.sql",
	}
	// ch := make(chan int, Num)
	db, err := conf.NewMysqlDb()
	if err != nil {
		log.Fatal(err)
	}
	db.GetConnections()
	for i := 0; i < Num; i++ {
		// go func(i int) {
		// 	db.Insert("insert into test(name, age) values(?,?)", fmt.Sprintf("test%d", i), i)
		// 	ch <- 1
		// }(i)
		db.Insert("insert into test(name, age) values(?,?)", fmt.Sprintf("test%d", i), i)
		// ch <- 1
	}

	// for i := 0; i < Num; i++ {
	// 	<-ch
	// }
	log.Println("mysql5:", time.Since(start).Seconds())
	wg.Done()
}
