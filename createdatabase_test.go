package gomysql

import "testing"

func TestCreate(t *testing.T) {
	schema := `
	CREATE TABLE person (
		first_name varchar(30) not null default '',
		last_name varchar(30) not null default '',
		email varchar(50) not null default ''
	);
	
	CREATE TABLE place (
		country varchar(50) not null default '',
		city varchar(50) not null default '',
		telcode int not null default 86
	)
		`
	conf := Sqlconfig{
		UserName:        "test",
		Password:        "123456",
		Port:            3306,
		Host:            "192.168.101.4",
		MultiStatements: true,
	}
	db, err := conf.NewDb()
	if err != nil {
		t.Fatal(err)
	}

	stat := db.Stats()
	t.Log(stat)
	ndb, err := db.Use("test")
	if err != nil {
		t.Fatal(err)
	}
	// 建表

	_, err = ndb.Query(schema)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success")
}
