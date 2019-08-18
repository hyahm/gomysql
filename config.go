package gomysql

import (
	"context"
	"database/sql"
	"fmt"
)

type Sqlconfig struct {
	UserName string
	Password string
	Host     string
	Port     int
	DbName   string
}
// 如果tag 是空的, 那么默认dbname
func SaveConf(tag string, c *Sqlconfig) error {
	//判断是否是空map
	if client == nil {
		client = make(map[string]*Db, 0)
	}
	connstring := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		c.UserName, c.Password, c.Host, c.Port, c.DbName,
	)

	// 验证配置有效性, 并写入连接

	conn, err := sql.Open("mysql", connstring)
	if err != nil {
		return err
	}

	if err = conn.Ping(); err != nil {
		return err
	}
	if tag == "" {
		tag = c.DbName
	}
	d := &Db {
		conn: conn,
		conf: connstring,
		key:tag,
		Ctx:context.Background(),
	}

	// 写入key
	//保存配置
	client[tag] = d
	return nil
}
