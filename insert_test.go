package gomysql

import (
	"testing"

	"github.com/hyahm/golog"
)

func TestInsert(t *testing.T) {

	defer golog.Sync()
	conf := Sqlconfig{
		UserName:        "test",
		Password:        "123456",
		Port:            3306,
		DbName:          "test",
		Host:            "192.168.101.4",
		MultiStatements: true,
	}
	golog.Info(11111111)
	db, err := conf.NewDb()
	if err != nil {
		t.Fatal(err)
	}
	ps := &Person{
		FirstName: "cander",
		LastName:  "biao",
		Email:     "aaaaa@eaml.com",
	}

	err = db.InsertWithoutID(ps, "insert into person values(?)")
	if err != nil {
		t.Fatal(err)
	}
}
