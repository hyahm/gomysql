package gomysql

import "database/sql"

type Rows struct {
	*sql.Rows
}

type Row struct {
	*sql.Row
	err error
}

func (r *Rows) Close() {
	r.Close()
}

func (r *Row) Close() {
	r.Scan()
	r.Close()
}

func (r *Row) Scanf(dest ...interface{}) error {
	// 请使用 Scanf 代替 Scan
	if r.err != nil {
		return r.err
	}
	return r.Scan(dest...)
}
