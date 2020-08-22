# gomysql
mysql 只是简单封装
 - [x] 支持高并发
 - [x] 支持更新和删除失败的日志记录
 - [x] 支持驱动自带的连接池，
 - [x] 避免连接过多导致的失败
 - [x] 支持in的操作

bench
```go
package main

import (
	"fmt"
	"github.com/hyahm/gomysql"
)

const Num = 1000

func main() {
        wg := &sync.WaitGroup{}
	wg.Add(1)
	go Insert8(wg)
	// go Insert5(wg)
	wg.Wait()
}

func Insert8(wg *sync.WaitGroup) {

	start := time.Now()
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
	ch := make(chan int, Num)
	db, err := conf.NewDb()
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
	
	rows,err := db.GetRowsIn("select id from test where age in (?)", gomysql.InArgs([]string{"1","2","3","4","5"}).ToInArgs())
	if err != nil {
		os.Exit(1)
	}
	for rows.Next() {
		var id int64
	  	rows.Scan(&id)
	}
	log.Println("mysql8:", time.Since(start).Seconds())
	wg.Done()
}

```




