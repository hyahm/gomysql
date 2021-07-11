package gomysql

import (
	"reflect"
	"testing"

	"github.com/hyahm/golog"
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
