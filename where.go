package gomysql

type Condition map[string]interface{}
type AndWhere map[string]interface{}
type OrWhere map[string]interface{}

type Where struct {
	Or  map[string]interface{}
	And map[string]interface{}
}

func NewWhere() *Where {
	return &Where{
		Or:  make(map[string]interface{}),
		And: make(map[string]interface{}),
	}
}
