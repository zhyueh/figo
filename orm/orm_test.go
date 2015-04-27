package orm

import (
	"model"
	"testing"
)

type user struct {
	orm.Model
	Id       int    `orm:"auto"`
	Name     string `orm:"varchar(16)" name:"name"`
	Age      int    `orm:"int" name:"age"`
	Datetime string `orm:"datetime" name:"datetime" empty:"ignore"`
}

func TestAll() {
	orm := getOrm()
	u := new(user)

	fmt.Println(orm.Fork().All(u))
}

func getOrm() *orm {
	orm, _ := NewOrm("mysql", "127.0.0.1", "root", "root", 3306)
	return orm
}
