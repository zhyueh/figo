package orm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DLMysql struct {
	db *sql.DB
	tx *sql.Tx
}

func NewMysqlLayer(driver, connString string) (*DLMysql, error) {
	db, err := sql.Open(driver, connString)
	if err != nil {
		return nil, err
	}

	p := new(DLMysql)
	p.db = db

	return p, nil
}

func myScan(rows *sql.Rows) DbRow {
	r := DbRow{}

	cols, _ := rows.Columns()
	c := len(cols)
	vals := make([]interface{}, c)
	valPtrs := make([]interface{}, c)

	for i := range cols {
		valPtrs[i] = &vals[i]
	}

	rows.Scan(valPtrs...)

	for i := range cols {
		if val, ok := vals[i].([]byte); ok {
			r[cols[i]] = string(val)
		} else {
			r[cols[i]] = vals[i]
		}
	}

	return r
}

func (this *DLMysql) Close() {
	this.db.Close()
}

//func (this *DLMysql) Transaction(fn func(*DLMysql) error) error {
//	if db, ok := p.Db.(*sql.DB); ok {
//		if tx, err := db.Begin(); err != nil {
//			return err
//		} else {
//			if err = fn(WrapDLMysql(tx)); err != nil {
//				tx.Rollback()
//				return err
//			} else {
//				tx.Commit()
//			}
//		}
//	}
//	return nil
//}

func (this *DLMysql) TBegin() error {
	if tx, err := this.db.Begin(); err != nil {
		return err
	} else {
		this.tx = tx
	}
	return nil
}

func (this *DLMysql) TCommit() error {
	return this.tx.Commit()
}

func (this *DLMysql) TRollback() error {
	return this.tx.Rollback()
}

func (this *DLMysql) TExec(sql string, args ...interface{}) (sql.Result, error) {
	return this.tx.Exec(sql, args...)
}

func (this *DLMysql) TQuery(sql string, args ...interface{}) ([]DbRow, error) {
	rows, err := this.tx.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := make([]DbRow, 0)

	for rows.Next() {
		r = append(r, myScan(rows))
	}

	return r, nil
}

func (this *DLMysql) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return this.db.Exec(sql, args...)
}

func (this *DLMysql) Insert(sql string, args ...interface{}) (int64, error) {
	res, err := this.db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	r, _ := res.LastInsertId()
	return r, nil
}

func (this *DLMysql) Update(sql string, args ...interface{}) (int64, error) {
	res, err := this.db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	r, _ := res.RowsAffected()
	return r, nil
}

func (this *DLMysql) Delete(sql string, args ...interface{}) (int64, error) {

	res, err := this.db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	r, _ := res.RowsAffected()
	return r, nil
}

func (this *DLMysql) One(sql string, args ...interface{}) (DbRow, error) {
	rows, err := this.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()

	return myScan(rows), nil
}

func (this *DLMysql) All(sql string, args ...interface{}) ([]DbRow, error) {
	rows, err := this.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := make([]DbRow, 0)

	for rows.Next() {
		r = append(r, myScan(rows))
	}

	return r, nil
}
