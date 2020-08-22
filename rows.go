package gomysql

import "database/sql"

type Row struct {
	*sql.Row
	err error
}

func (r *Row) Scanf(dest ...interface{}) error {
	// 请使用 Scanf 代替 Scan, 多封装了一层error
	if r.err != nil {
		return r.err
	}
	return r.Scan(dest...)
}
