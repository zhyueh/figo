package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/zhyueh/figo/toolkit"
	"reflect"
)

type Orm struct {
	db     dlInterface
	dbType string
	qb     *MySQLQB

	host     string
	user     string
	password string
	name     string
	port     int
}

func NewOrm(dbType, host, user, password, name string, port int) (*Orm, error) {
	re := new(Orm)
	re.dbType = dbType
	re.host = host
	re.user = user
	re.password = password
	re.name = name
	re.port = port

	re.qb = NewMySQLQB()

	if db, err := NewDL(dbType, host, user, password, name, port); err != nil {
		return nil, err
	} else {
		re.db = db
	}

	return re, nil
}

func (this *Orm) Fork() *Orm {
	re := new(Orm)
	re.db = this.db
	re.dbType = this.dbType
	re.qb = NewMySQLQB()

	return re

}

/*
 * need close connection when using clone
 */
func (this *Orm) Clone() (*Orm, error) {
	newOrm, err := NewOrm(
		this.dbType,
		this.host,
		this.user,
		this.password,
		this.name,
		this.port,
	)
	if err != nil {
		return nil, err
	}
	return newOrm, nil
}

func (this *Orm) Transaction() (*Orm, error) {
	//new db connection for transaction
	newOrm, err := this.Clone()
	if err != nil {
		return nil, err
	}

	transactionError := newOrm.db.TBegin()
	if transactionError != nil {
		return nil, transactionError
	}

	return newOrm, nil
}

func (this *Orm) TCommit() error {
	return this.db.TCommit()
}

func (this *Orm) TRollback() error {
	return this.db.TRollback()
}

func (this *Orm) TExec(sql string, args ...interface{}) (sql.Result, error) {
	return this.db.TExec(sql, args...)
}

func (this *Orm) TQuery(sql string, args ...interface{}) ([]DbRow, error) {
	return this.db.TQuery(sql, args...)
}

func (this *Orm) Close() {
	this.db.Close()
}

func (this *Orm) Where(cond string, i ...interface{}) *Orm {
	this.qb.Where(cond, i...)
	return this
}

func (this *Orm) WhereIn(cond string, i ...interface{}) *Orm {
	this.qb.WhereIn(cond, i...)
	return this
}

func (this *Orm) And(cond string, i interface{}) *Orm {
	this.qb.And(cond, i)
	return this
}

func (this *Orm) Or(cond string, i interface{}) *Orm {
	this.qb.Or(cond, i)
	return this
}

func (this *Orm) Order(order string) *Orm {
	this.qb.Order(order)
	return this
}

func (this *Orm) Page(index, num int) *Orm {
	this.qb.Page(index, num)
	return this
}

func (this *Orm) QueryRaw(o interface{}, sql string, args ...interface{}) (error, bool) {
	dbrow, err := this.db.One(sql, args...)
	if err != nil {
		return err, false
	}
	DbRowToModel(dbrow, o)
	return nil, len(dbrow) != 0
}

func (this *Orm) Find(o ModelInterface) (error, bool) {
	defer this.qb.Reset()
	field, val, exists := GetIdFieldValue(o)
	if exists {
		this.qb.Where(fmt.Sprintf("`%s`=?", field), val)
		this.qb.Table(ModelTableName(o))
		sql, args := this.qb.Select()
		//fmt.Println(sql)
		dbrow, err := this.db.One(sql, args...)
		if err != nil {
			return err, false
		}
		DbRowToModel(dbrow, o)
		return nil, len(dbrow) != 0
	} else {
		return errors.New("no auto fields"), false
	}
}

func (this *Orm) One(o ModelInterface) (error, bool) {
	defer this.qb.Reset()
	this.qb.Table(ModelTableName(o))
	sql, args := this.qb.Select()
	//fmt.Println(sql)
	dbrow, err := this.db.One(sql, args...)
	if err != nil {
		return err, false
	}
	DbRowToModel(dbrow, o)
	return nil, len(dbrow) != 0
}

func (this *Orm) All(o ModelInterface) (error, []interface{}) {
	defer this.qb.Reset()
	//if len(os) == 0 {
	//	return errors.New("no data list when calling All"), 0
	//}
	//o := os[0]
	re := make([]interface{}, 0)
	val := reflect.ValueOf(o).Elem()
	modelType := val.Type()

	this.qb.Table(ModelTableName(o))
	sql, args := this.qb.Select()
	//fmt.Println(sql)
	dbrows, err := this.db.All(sql, args...)
	if err != nil {
		return err, re
	}
	if len(dbrows) == 0 {
		return nil, re
	}
	for _, dbrow := range dbrows {
		newo := reflect.New(modelType).Interface()
		DbRowToModel(dbrow, newo)
		re = append(re, newo)
	}

	return nil, re
}

func (this *Orm) Count(o ModelInterface) (int, error) {
	defer this.qb.Reset()
	this.qb.Table(ModelTableName(o))
	sql, args := this.qb.Count()
	dbrows, err := this.db.All(sql, args...)
	//fmt.Println(sql)
	//fmt.Println(dbrows)
	if err != nil || len(dbrows) != 1 {
		return 0, err
	} else {
		if val, exists := dbrows[0]["num"]; exists {
			return toolkit.ConvertToInt(val), nil
		}
		return 0, nil
	}
}

func (this *Orm) Save(o ModelInterface) error {
	defer this.qb.Reset()
	fields, values := GetSaveModelFieldValues(o)
	this.qb.Fields(fields)
	this.qb.Values(values)
	//fmt.Println("model table name", ModelTableName(o))
	this.qb.Table(ModelTableName(o))

	if NeedInsertModel(o) {
		sql, args := this.qb.Insert()
		//fmt.Println(sql)
		id, err := this.db.Insert(sql, args...)
		if err != nil {
			return err
		}
		ModelUpdateId(o, id)
	} else {
		kFields, kValues, _ := GetKeyFieldValues(o)
		//fmt.Println(kFields, kValues)
		if len(kFields) == 0 {
			return errors.New("no id define")
		}
		this.qb.KeyFields(kFields)
		this.qb.KeyValues(kValues)
		sql, args := this.qb.InsertIgnore()
		//fmt.Println(sql, args)
		_, err := this.db.Update(sql, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Orm) QueryOne(sql string, args ...interface{}) (DbRow, error) {
	return this.db.One(sql, args...)
}

func (this *Orm) Query(sql string, args ...interface{}) ([]DbRow, error) {
	return this.db.All(sql, args...)
}

func (this *Orm) Execute(sql string, args ...interface{}) error {
	_, err := this.db.Exec(sql, args...)
	return err
}

func (this *Orm) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return this.db.Exec(sql, args...)
}
