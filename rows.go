package gomysql

import "database/sql"

type Row struct {
	Row *sql.Row
	err error
}

type Rows struct {
	Rows *sql.Rows
	err  error
}

func (r *Row) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	return r.Row.Scan(dest...)
}

func (rs *Rows) Scan(dest ...interface{}) error {
	return rs.Rows.Scan(dest...)
}

func (rs *Rows) Next() bool {
	return rs.Rows.Next()
}

func (rs *Rows) Close(dest ...interface{}) error {

	return rs.Rows.Close()
}

func (rs *Rows) Err() error {
	return rs.err
}

func (rs *Rows) NextResultSet() bool {
	return rs.Rows.NextResultSet()
}

func (rs *Rows) ColumnTypes() ([]*sql.ColumnType, error) {
	return rs.Rows.ColumnTypes()
}

func (rs *Rows) Columns() ([]string, error) {
	return rs.Rows.Columns()
}
