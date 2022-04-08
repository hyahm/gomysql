package gomysql

import (
	"testing"
)

type MeStruct struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type Person struct {
	ID        int64    `db:"id,default"`
	FirstName string   `db:"first_name,force"`
	LastName  string   `db:"last_name"`
	Email     string   `db:"email,force"`
	Me        MeStruct `db:"me"`
	Uids      []int64  `db:"uids"`
	TestJson  string   `db:"test"`
	Age       int      `json:"age" db:"age,counter"`
}

func TestInsert(t *testing.T) {

	conf := Sqlconfig{
		UserName:        "test",
		Password:        "123456",
		Port:            3306,
		DbName:          "test",
		Host:            "192.168.50.250",
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
		Uids:     []int64{5, 7, 23, 12, 90},
		Age:      1,
	}
	// pss := make([]*Person, 0)
	// for i := 0; i < 20; i++ {
	// 	pss = append(pss, ps)
	// }
	// $key  $value 是固定占位符
	// default: 如果设置了并且为零值， 那么为数据库的默认值
	// struct, 指针， 切片 默认值为 ""
	res := db.InsertInterfaceWithoutID(ps, "insert into person($key) values($value)")
	if res.Err != nil {
		t.Fatal(err)
	}
}
