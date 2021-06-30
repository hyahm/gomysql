package gomysql

import "testing"

func TestWhere(t *testing.T) {
	where := NewWhere()

	where.Or["id"] = 6
}
