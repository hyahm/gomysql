# gomysql
mysql 只是简单封装

# 初衷
- 抛弃各种orm
- 通过配置来驱动
- 建议将连接保存某个全局变量， 所有的地方都可以直接执行， 执行完成可以关闭也可以不关闭连接， 不关闭就是长链接， 
- 关闭就是短链接， 下次调用也能直接调用
- 增加调试sql


v0.0.2 版
- 删除tx的支持， 需要使用的话，通过Db.Begin() 自行生成
- 增加运行sql的调试信息， 可以打印运行的sql，方便找出sql错误
- 减少复杂调用
example.go
```
package main

import (
	"fmt"
	"github.com/hyahm/gomysql"
)

var (
	conf = &gomysql.Sqlconfig{
		Host: "127.0.0.1",
		Port: 3306,
		UserName: "zth",
		Password: "123456",
		DbName: "zth",
	}
)


func main() {
	db, err := conf.NewDb()
	if err != nil {
		panic(err)
	}
	var id int64
	db.OpenDebug()
	err = db.GetOne("select id from cmf_developer limit 1").Scan(&id)
	if err != nil {
	// todo
	}
	
	fmt.Println(db.PrintSql())
	fmt.Println(id)
}

```
out
```
select id from cmf_developer limit 1
1
```


