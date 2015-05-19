package orm

import (
	"github.com/zhyueh/figo/toolkit"
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
	Timestamp time.Time `orm:"datetime" readonly:"true" name:"timestamp" empty:"ignore"`
}

type UserDetail struct {
	Model
	UserId  int    `orm:"primary"`
	Address string `orm:"varchar(25)"`
}

func TestInsertIgnore(t *testing.T) {
	orm := getOrm()
	ud := new(UserDetail)
	ud.UserId = 1
	ud.Address = toolkit.RandomString(5)
	err := orm.Fork().Save(ud)
	if err != nil {
		t.Fatal("insert ignore error", err)
	}
}

func TestSaveAndFind(t *testing.T) {
	orm := getOrm()
	randomString := toolkit.RandomString(5)
	u := new(user)
	u.Name = randomString
	err := orm.Fork().Save(u)
	if err != nil {
		t.Fatal(err)
	}
	if u.Id < 1 {
		t.Fatal("new id ", u.Id)
	}

	newu := new(user)
	newu.Id = u.Id
	newerr, exists := orm.Fork().Find(newu)
	if newerr != nil {
		t.Fatal(newerr)
	}

	if !exists {
		t.Fatal("save but find ", newu.Id, u.Id)
	}

	if newu.Name != randomString {
		t.Fatal("find diff value")
	}
}

func TestCount(t *testing.T) {
	orm := getOrm()
	u := new(user)
	count, err := orm.Fork().Count(u)
	t.Log("count:", count)
	if err != nil {
		t.Fatal("count error", err)
	}
	if count < 1 {
		t.Fatal("count less than zero")
	}

	count, err = orm.Fork().Where("`id` < ?", 3).Count(u)
	if err != nil {
		t.Fatal("count where error", err)
	}
	if count != 2 {
		t.Fatal("count where is not one")
	}
	t.Log("count where:", count)

}

func TestPage(t *testing.T) {
	orm := getOrm()
	u := new(user)
	count, err := orm.Fork().Count(u)

	if count < 1 || err != nil {
		t.Fatal("can not test page because count error :", count, err)
	}

	page := 0
	num := 5
	total_page := count / num
	if count%num != 0 {
		total_page += 1
	}

	for ; page < total_page; page++ {
		page_err, users := orm.Fork().Page(page, num).All(u)
		if page_err != nil {
			t.Fatal("page error:", page, num, page_err)
		}
		for _, tmp := range users {
			t.Log(tmp.(*user))
		}
		t.Log(page, "end")
	}

}

func getOrm() *Orm {
	orm, _ := NewOrm("mysql", "127.0.0.1", "root", "root", "test", 3306)
	return orm
}
