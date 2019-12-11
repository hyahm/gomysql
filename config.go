package gomysql

import (
	"context"
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
func (s *Sqlconfig) NewDb() (*Db,error) {
	//判断是否是空map
	connstring := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		s.UserName, s.Password, s.Host, s.Port, s.DbName,
	)
	db := &Db{
		conf: connstring,
		Ctx: context.Background(),
	}

	return db.conndb()
}
