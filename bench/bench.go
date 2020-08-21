package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/gomysql"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	Insert8(wg)
	go Insert5(wg)
	wg.Wait()
}

const Num = 1000000

func Insert8(wg *sync.WaitGroup) {
	fmt.Println(111)
	wg.Add(1)
	start := time.Now()
	conf := &gomysql.Sqlconfig{
		Host:               "192.168.50.211",
		UserName:           "cander",
		Password:           "123456",
		DbName:             "test",
		Port:               3306,
		MaxOpenConns:       1000,
		MaxIdleConns:       1,
		WriteLogWhenFailed: true,
		LogFile:            ".failedlinux.sql",
	}
	ch := make(chan int, Num)
	db, err := conf.NewDb()
	if err != nil {
		golog.Error(err)
		os.Exit(1)
	}
	for i := 0; i < Num; i++ {
		go func() {
			_, err := db.Insert("insert into test(name, age) values(?,?)", fmt.Sprintf("test%d", i), i)
			if err != nil {
				golog.Error(err)
			}
			ch <- 1
		}()

	}

	for i := 0; i < Num; i++ {
		<-ch
	}
	log.Println("mysql8:", time.Since(start).Seconds())
	wg.Done()
}

func Insert5(wg *sync.WaitGroup) {
	wg.Add(1)
	start := time.Now()
	conf := &gomysql.Sqlconfig{
		Host:               "192.168.50.49",
		UserName:           "test",
		Password:           "123456",
		DbName:             "test",
		Port:               3306,
		MaxOpenConns:       1000,
		MaxIdleConns:       1,
		WriteLogWhenFailed: true,
		LogFile:            ".failedwindows.sql",
	}
	ch := make(chan int, Num)
	db, err := conf.NewDb()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < Num; i++ {
		go func() {
			db.Insert("insert into test(name, age) values(?,?)", fmt.Sprintf("test%d", i), i)
			ch <- 1
		}()

	}

	for i := 0; i < Num; i++ {
		<-ch
	}
	log.Println("mysql5:", time.Since(start).Seconds())
	wg.Done()
}
