package orm

import (
	"errors"
	"fmt"
)

type Orm struct {
	db     dlInterface
	dbType string
	qb     *MySQLQB
}

func NewOrm(dbType, host, user, password, name string, port int) (*Orm, error) {
	re := new(Orm)

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

func (this *Orm) Close() {
	this.db.Close()
}

func (this *Orm) Where(cond string, i interface{}) *Orm {
	this.qb.Where(cond, i)
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

func (this *Orm) Find(o ModelInterface) (error, bool) {
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

func (this *Orm) One(o ModelInterface) error {
	return nil
}

func (this *Orm) Save(o ModelInterface) error {
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
		idName, idVal, idExists := GetIdFieldValue(o)
		if !idExists {
			return errors.New("no id define")
		}
		this.qb.Where(
			fmt.Sprintf("`%s`=?", idName),
			idVal)
		sql, args := this.qb.Update()
		fmt.Println(sql)
		_, err := this.db.Update(sql, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Orm) Query(sql string, args ...interface{}) ([]DbRow, error) {
	return this.db.All(sql, args...)
}
