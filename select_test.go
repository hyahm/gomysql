package gomysql

import (
	"testing"

	"github.com/hyahm/golog"
)

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

func TestSelect(t *testing.T) {
	defer golog.Sync()
	conf := Sqlconfig{
		UserName:        "test",
		Password:        "123456",
		Port:            3306,
		DbName:          "test",
		Host:            "192.168.101.4",
		MultiStatements: true,
	}
	db, err := conf.NewDb()
	if err != nil {
		t.Fatal(err)
	}
	persons := make([]Person, 0)
	err = db.Select(&persons, "select * from person")
	t.Log(err)
	t.Log(persons)
	// 建表

}
