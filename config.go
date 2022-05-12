package gomysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
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
	MaxAllowedPacket        uint64 // insert 插入大量数据的时候会用到 或者用insertmany
	Collation               string
	MaxOpenConns            int           // 请设置小于等于mysql 的max_connections值， 建议等于max_connections
	MaxIdleConns            int           // 如果设置了 MaxOpenConns， 那么此直将等于 MaxOpenConns
	ConnMaxLifetime         time.Duration // 连接池设置
	WriteLogWhenFailed      bool
	LogFile                 string
	Debug                   bool // 打印sql
}

func (s *Sqlconfig) GetMysqlDataSource() string {
	s.setDefaultConfig()
	//判断是否是空map
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&clientFoundRows=%t&allowCleartextPasswords=%t&interpolateParams=%t&columnsWithAlias=%t&multiStatements=%t&parseTime=%t&tls=%t&readTimeout=%s&timeout=%s&allowOldPasswords=%t&loc=%s&maxAllowedPacket=%d&collation=%s&writeTimeout=%s",
		s.UserName, s.Password, s.Host, s.Port, s.DbName, s.Charset, s.ClientFoundRows, s.AllowCleartextPasswords, s.InterpolateParams, s.ColumnsWithAlias, s.MultiStatements, s.ParseTime, s.TLS, s.ReadTimeout, s.Timeout, s.AllowCleartextPasswords, s.Loc, s.MaxAllowedPacket, s.Collation, s.WriteTimeout,
	)
}

func (s *Sqlconfig) GetPostgreDataSource() string {
	s.setpgDefaultConfig()
	//判断是否是空map
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		s.UserName, s.Password, s.Host, s.Port, s.DbName,
	)
}

// 如果tag 是空的, 那么默认dbname
func (s *Sqlconfig) NewMysqlDb() (*Db, error) {
	return s.conndb(s.GetMysqlDataSource())
}

func (s *Sqlconfig) NewPGPool() (*PGConn, error) {
	return s.connPg(s.GetPostgreDataSource())
}

// 不存在就创建database
func (s *Sqlconfig) CreateDB(name string) (*Db, error) {
	db, err := s.conndb(s.GetMysqlDataSource())
	if err != nil {
		return nil, err
	}
	newdb, err := db.Use(name)
	if err != nil {
		return nil, err
	}
	db.Close()

	return newdb, err
}

func (s *Sqlconfig) connPg(conf string) (*PGConn, error) {
	conn, err := pgxpool.Connect(context.Background(), s.GetPostgreDataSource())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err = conn.Ping(context.Background()); err != nil {
		conn.Close()
		return nil, err
	}

	db := &PGConn{
		conn,
		conf,
		context.Background(),
		s,
		nil,
		s.Debug,
	}

	if s.ReadTimeout == 0 {
		s.ReadTimeout = time.Second * 30
	}

	// 防止开始就有很多连接，导致
	// ch = make(chan struct{}, db.maxConn)

	return db, nil
}

func (s *Sqlconfig) conndb(conf string) (*Db, error) {

	conn, err := sql.Open("mysql", conf)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	db := &Db{
		conn,
		conf,
		context.Background(),
		s,
		nil,
		nil,
		0,
		0,
		s.Debug,
	}

	if s.ReadTimeout == 0 {
		s.ReadTimeout = time.Second * 30
	}
	db.SetMaxIdleConns(s.MaxIdleConns)
	if s.MaxOpenConns > 0 {
		db.maxConn = s.MaxOpenConns
	} else {
		db.maxConn = 100
	}
	db.maxpacket = s.MaxAllowedPacket
	if s.MaxAllowedPacket == 0 {
		db.maxpacket = 4 * 1024 * 1024 // 默认4m
	}
	// 防止开始就有很多连接，导致
	// ch = make(chan struct{}, db.maxConn)

	db.SetMaxOpenConns(s.MaxOpenConns)
	db.SetConnMaxLifetime(s.ConnMaxLifetime)
	if db.sc.WriteLogWhenFailed {
		db.mu = &sync.RWMutex{}
		if db.sc.LogFile == "" {
			db.sc.LogFile = ".failed.sql"
		}
		var err error
		db.f, err = os.OpenFile(db.sc.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	return db, nil
}

func (s *Sqlconfig) setDefaultConfig() {
	if s.Charset == "" {
		s.Charset = "utf8mb4"
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

func (s *Sqlconfig) setpgDefaultConfig() {
	if s.UserName == "" {
		s.UserName = "postgres"
	}
	if s.Host == "" {
		s.Host = "127.0.0.1"
	}
	if s.Port == 0 {
		s.Port = 5432
	}
	if s.DbName == "" {
		s.DbName = "postgres"
	}
}
