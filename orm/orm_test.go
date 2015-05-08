package orm

import (
	"testing"
	"time"
)

type user struct {
	Model
	Id   int    `orm:"auto"`
	Name string `orm:"varchar(16)" name:"name"`
	Age  int    `orm:"int" name:"age" readonly:"1"`
	//Datetime  time.Time `orm:"datetime" name:"datetime" empty:"ignore"`
	//Date      time.Time `orm:"datetime" name:"date" empty:"ignore"`
	Timestamp time.Time `orm:"datetime" name:"timestamp" empty:"ignore"`
}

func TestAll(t *testing.T) {
	orm := getOrm()
	//u := new(user)
	//t.Log(orm.Fork().All(u))
	u := new(user)
	u.Id = 4
	orm.Fork().Find(u)
	t.Fatal(u)
}

func getOrm() *Orm {
	orm, _ := NewOrm("mysql", "127.0.0.1", "root", "root", "test", 3306)
	return orm
}
