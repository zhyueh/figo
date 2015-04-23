package orm

import (
	"errors"
	"fmt"
)

type DbRow map[string]interface{}

type dlInterface interface {
	Close()
	Exec(string, ...interface{}) error
	Insert(string, ...interface{}) (int64, error)
	Update(string, ...interface{}) (int64, error)
	Delete(string, ...interface{}) (int64, error)
	One(string, ...interface{}) (DbRow, error)
	All(string, ...interface{}) ([]DbRow, error)
}

func NewDL(dbType, host, user, password, db string, port int) (dlInterface, error) {
	if dbType == "mysql" {
		connectString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, password, host, port, db)
		return NewMysqlLayer(dbType, connectString)
	}
	return nil, errors.New("un support db layer type")
}
