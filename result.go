package gomysql

import "errors"

// 返回的结果
type Result struct {
	Err           error
	Sql           string // sql并不是真实的执行的sql， 仅供参考
	LastInsertId  int64
	LastInsertIds []int64
	// RowsAffected returns the number of rows affected by an
	// update, insert, or delete. Not every database or database
	// driver may support this.
	RowsAffected int64
}

var ErrNotSupport = errors.New("dest type not support")
