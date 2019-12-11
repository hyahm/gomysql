# gomysql
mysql 只是简单封装

# 初衷
- 抛弃各种orm
- 通过配置来驱动
- 建议将连接保存某个全局变量， 所有的地方都可以直接执行， 执行完成可以关闭也可以不关闭连接， 不关闭就是长链接， 
- 关闭就是短链接， 下次调用也能直接调用
- 增加调试sql

# 使用
```
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
  	// 保存配置到配置文件, 保存再map 中
	gomysql.SaveConf("x7", Conf) 
  	// 连接指定key的配置文件, 保存再map 中
	x7 ,err = gomysql.ConnDB("x7")
	if err != nil {
		panic("connect mysql error")
	}
	// 上面的连接使用后,  后面的内容在其后执行的代码任何地方都能调同
	rows,err := x7.GetRows("select username,password from user")
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
	x7.Close()  // 关闭连接，  如果不关闭就是长连接， 注意， 就算关闭了也可以其他的地方调用sql 语句，  会根据保存的配置从新连接
}
```

# 安装
```
go get github.com/hyahm/gomysql
```
推荐直接复制sql.go 文件的内容到项目,  自己定制


v0.0.2 版
- 删除tx的支持， 需要使用的花，通过Db.Begin() 自行生成
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
		panic(err)
	}
	fmt.Println(db.PrintSql())
	fmt.Println(id)
}

```


v0.0.2 与v0.0.3 的区别
因为源码sql的问题导致, 看起来 0.0.2更简单， 但是如果重联出错，会panic
v0.0.2
```
err := db.GetOne("select id from cmf_developer limit 1").Scan(&id)
```
v0.0.3
```
row, err := db.GetOne("select id from cmf_developer limit 1")
if err != nil {
// todo
}
err = row.Scan(&id)
```



