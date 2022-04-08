package gomysql

import (
	"fmt"
)

type Curder interface {
	Read(model interface{}) Result
	Update(interface{}, string) Result
	Create(model interface{}) Result
	Delete(interface{}, string) Result
}

type Actuator struct {
	Table string
	Db    *Db
}

func (db *Db) NewCurder(table string) Curder {
	return &Actuator{
		Table: table,
		Db:    db,
	}
}

func (actuator *Actuator) Read(model interface{}) Result {
	return actuator.Db.Select(model, fmt.Sprintf("select * from %s", actuator.Table))
}

func (actuator *Actuator) Update(model interface{}, where string) Result {
	return actuator.Db.UpdateInterface(model, "update %s set $set where %s", actuator.Table, where)
}

func (actuator *Actuator) Create(model interface{}) Result {
	res := actuator.Db.InsertInterfaceWithID(model, fmt.Sprintf("insert into %s($key) values($value)", actuator.Table))
	return res
}

func (actuator *Actuator) Delete(model interface{}, where string) Result {
	return actuator.Db.Select(model, fmt.Sprintf("delete from %s where %s", actuator.Table, where))
}
