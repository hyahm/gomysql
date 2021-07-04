package gomysql

import (
	"testing"

	"github.com/hyahm/golog"
)

func TestUpdate(t *testing.T) {
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
	ps := &Person{
		FirstName: "what is it",
		LastName:  "hyahm.com",
		Email:     "aaaaa@eaml.com",
		Me: MeStruct{
			X: 10,
			Y: 20,
			Z: 30,
		},
	}

	// $key  $value 是固定占位符
	// omitempty: 如果为空， 那么为数据库的默认值
	// struct, 指针， 切片 默认值为 ""
	// $set
	golog.Info("start update")
	_, err = db.UpdateInterface(ps, "update person set $set where id=?", 1)
	if err != nil {
		golog.Error(err)
		t.Fatal(err)
	}
}
