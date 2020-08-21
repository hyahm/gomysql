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
	Tls                     bool
	ReadTimeout             string
	Timeout                 string
	AllowOldPasswords       bool
	Charset                 string
	Loc                     string
	MaxAllowedPacket        uint64
	Collation               string
	MaxOpenConns            int
	MaxIdleConns            int
	ConnMaxLifetime         time.Duration
	WriteLogWhenFailed      bool
	LogFile                 string
}

// 如果tag 是空的, 那么默认dbname
func (s *Sqlconfig) NewDb() (*Db, error) {
	s.setDefaultConfig()
	//判断是否是空map
	connstring := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&clientFoundRows=%t&allowCleartextPasswords=%t&interpolateParams=%t&columnsWithAlias=%t&multiStatements=%t&parseTime=%t&tls=%t&readTimeout=%s&timeout=%s&allowOldPasswords=%t&loc=%s&maxAllowedPacket=%d&collation=%s",
		s.UserName, s.Password, s.Host, s.Port, s.DbName, s.Charset, s.ClientFoundRows, s.AllowCleartextPasswords, s.InterpolateParams, s.ColumnsWithAlias, s.MultiStatements, s.ParseTime, s.Tls, s.ReadTimeout, s.Timeout, s.AllowCleartextPasswords, s.Loc, s.MaxAllowedPacket, s.Collation,
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
	if s.Timeout == "" {
		s.Timeout = "1s"
	}
	if s.ReadTimeout == "" {
		s.ReadTimeout = "5s"
	}

}
