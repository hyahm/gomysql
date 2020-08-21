package gomysql

import "database/sql"

type Rows struct {
	*sql.Rows
}

type Row struct {
	*sql.Row
}

func (r *Rows) Close() {
	r.Close()
	<-ch
}

func (r *Row) Close() {
	r.Scan()
	r.Close()
}

func (r *Row) Scanf(dest ...interface{}) error {
	// 请使用 Scanf 代替 Scan
	defer func() {
		<-ch
	}()
	return r.Scan(dest...)
}
