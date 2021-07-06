package gomysql

import (
	"testing"

	"github.com/hyahm/golog"
)

type MeStruct struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type Person struct {
	ID        int64    `db:"id"`
	FirstName string   `db:"first_name,force"`
	LastName  string   `db:"last_name"`
	Email     string   `db:"email,omitempty,force"`
	Me        MeStruct `db:"me"`
	TestJson  string   `db:"test"`
}

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
	db, err := conf.NewDb()
	if err != nil {
		t.Fatal(err)
	}
	ps := &Person{
		FirstName: "cander",
		LastName:  "biao",
		Email:     "aaaaa@eaml.com",
		Me: MeStruct{
			X: 10,
			Y: 20,
			Z: 30,
		},
		TestJson: "testaaaa",
	}
	pss := make([]*Person, 0)
	for i := 0; i < 20; i++ {
		pss = append(pss, ps)
	}
	// $key  $value 是固定占位符
	// omitempty: 如果为空， 那么为数据库的默认值
	// struct, 指针， 切片 默认值为 ""
	err = db.InsertInterfaceWithoutID(pss, "insert into person($key) values($value)")
	if err != nil {
		golog.Error(err)
		t.Fatal(err)
	}
}
