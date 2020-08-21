package gomysql

import (
	"context"
	"fmt"
	"time"
)

type Sqlconfig struct {
	UserName                string
	Password                string
	Host                    string
	Port                    int
	DbName                  string
	ClientFoundRows         bool
	AllowCleartextPasswords bool
	InterpolateParams       bool
	ColumnsWithAlias        bool
	MultiStatements         bool
	ParseTime               bool
	TLS                     bool
	ReadTimeout             time.Duration
	Timeout                 time.Duration
	WriteTimeout            time.Duration
	AllowOldPasswords       bool
	Charset                 string
	Loc                     string
	MaxAllowedPacket        uint64
	Collation               string
	MaxOpenConns            int // 请设置小于mysql 的max_connections值
	MaxIdleConns            int
	ConnMaxLifetime         time.Duration
	WriteLogWhenFailed      bool
	LogFile                 string
}

// 如果tag 是空的, 那么默认dbname
func (s *Sqlconfig) NewDb() (*Db, error) {
	s.setDefaultConfig()
	//判断是否是空map
	connstring := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&clientFoundRows=%t&allowCleartextPasswords=%t&interpolateParams=%t&columnsWithAlias=%t&multiStatements=%t&parseTime=%t&tls=%t&readTimeout=%s&timeout=%s&allowOldPasswords=%t&loc=%s&maxAllowedPacket=%d&collation=%s&writeTimeout=%s",
		s.UserName, s.Password, s.Host, s.Port, s.DbName, s.Charset, s.ClientFoundRows, s.AllowCleartextPasswords, s.InterpolateParams, s.ColumnsWithAlias, s.MultiStatements, s.ParseTime, s.TLS, s.ReadTimeout, s.Timeout, s.AllowCleartextPasswords, s.Loc, s.MaxAllowedPacket, s.Collation, s.WriteTimeout,
	)
	db := &Db{
		conf: connstring,
		Ctx:  context.Background(),
		sc:   s,
	}

	return db.conndb()
}

func (s *Sqlconfig) setDefaultConfig() {
	if s.Charset == "" {
		s.Charset = "utf8"
	}
	if s.Collation == "" {
		s.Collation = "utf8mb4_general_ci"
	}
	if s.Loc == "" {
		s.Loc = "UTC"
	}
	if s.Timeout == 0 {
		s.Timeout = time.Second * 2
	}
	if s.ReadTimeout == 0 {
		s.ReadTimeout = time.Second * 20
	}

}
