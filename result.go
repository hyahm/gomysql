package gomysql

// 返回的结果
type Result struct {
	Err           error
	Sql           string
	LastInsertId  int64
	LastInsertIds []int64
	// RowsAffected returns the number of rows affected by an
	// update, insert, or delete. Not every database or database
	// driver may support this.
	RowsAffected int64
}
