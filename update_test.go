package gomysql

import (
	"reflect"
	"testing"
)

func TestInt(t *testing.T) {
	var i1 int = 0
	var i2 int16 = 0
	var i3 int8 = 0
	var i4 int32 = 0
	var i5 int64 = 0
	var i6 float32 = 0
	var i7 float64 = 0
	t.Log(reflect.ValueOf(i1).Int() == 0)
	t.Log(reflect.ValueOf(i2).Int() == 0)
	t.Log(reflect.ValueOf(i3).Int() == 0)
	t.Log(reflect.ValueOf(i4).Int() == 0)
	t.Log(reflect.ValueOf(i5).Int() == 0)
	t.Log(reflect.ValueOf(i6).Float() == 0)
	t.Log(reflect.ValueOf(i7).Float() == 0)

}

func TestStruct(t *testing.T) {
	p := Person{}
	t.Log(reflect.DeepEqual(p, Person{}))
}

func TestUpdate(t *testing.T) {
	conf := Sqlconfig{
		UserName:        "test",
		Password:        "123456",
		Port:            3306,
		DbName:          "test",
		Host:            "192.168.50.250",
		MultiStatements: true,
	}
	db, err := conf.NewMysqlDb()
	if err != nil {
		t.Fatal(err)
	}
	ps := &Person{
		FirstName: "what is it",
		LastName:  "hyahm.com",
		Email:     "aaaaa@eaml.com",
		// Me: MeStruct{
		// 	X: 10,
		// 	Y: 20,
		// 	Z: 30,
		// },
		Uids: []int64{1},
		Age:  1,
	}

	// $key  $value 是固定占位符
	// omitempty: 如果为空， 那么为数据库的默认值
	// struct, 指针， 切片 默认值为 ""
	// $set
	res := db.UpdateInterface(ps, "update person set $set where id=?", 4)
	if res.Err != nil {
		t.Fatal(err)
	}

}
