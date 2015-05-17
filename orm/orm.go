package orm

import (
	"errors"
	"fmt"
	"github.com/zhyueh/figo/toolkit"
	"reflect"
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

func (this *Orm) Where(cond string, i ...interface{}) *Orm {
	this.qb.Where(cond, i...)
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

func (this *Orm) One(o ModelInterface) (error, bool) {
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
	//if len(os) == 0 {
	//	return errors.New("no data list when calling All"), 0
	//}
	//o := os[0]
	re := make([]interface{}, 0)
	val := reflect.ValueOf(o).Elem()
	modelType := val.Type()

	this.qb.Table(ModelTableName(o))
	sql, args := this.qb.Select()
	fmt.Println(sql)
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
		//fmt.Println(sql)
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

func (this *Orm) Execute(sql string, args ...interface{}) error {
	return this.db.Exec(sql, args...)
}
