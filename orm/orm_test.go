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

type UserDetailEx struct {
	UserDetail `orm:"extend"`
	Abc        int `orm:"int"`
}

type UserDetailExEx struct {
	UserDetailEx `orm:"extend"`
}

func TestDbRowToModelEx(t *testing.T) {
	orm := getOrm()
	sql := "select * , 12 as abc from user_detail limit 1"
	if dbrows, err := orm.Query(sql); err == nil {
		for _, row := range dbrows {
			t.Log(row)
			ex := new(UserDetailExEx)
			DbRowToModelEx(row, ex)
			if ex.UserId == 0 {
				t.Fatal("dbrow to model ex not work")
			}
			t.Log(ex)
		}
	} else {
		t.Fatal(err)
	}

}

func FatalError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransaction(t *testing.T) {
	orm := getOrm()
	u := new(user)
	user_count, uc_error := orm.Fork().Count(u)
	t.Log("user count begin", user_count)
	if uc_error != nil {
		t.Fatal(uc_error)
	}

	if user_count == 0 {
		t.Fatal("can not test transction for no record in user")
	}

	defer orm.Close()
	tx, err := orm.Transaction()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Close()
	sqlResult, txError := tx.TExec("delete from user limit 1")
	FatalError(t, txError)
	deleted, _ := sqlResult.RowsAffected()
	if int(deleted) != 1 {
		t.Fatal("can not delete in transaction?")
	}
	t.Log("transction can delete nums", deleted)

	tx.TRollback()
	user_count_rollback, uc_error_rollback := orm.Fork().Count(u)
	FatalError(t, uc_error_rollback)
	if user_count_rollback != user_count {
		t.Fatal("can not rollback")
	}
	t.Log("user nums after rollback", user_count_rollback)

	//new transaction after rollback or commit
	tx, _ = orm.Transaction()
	sqlResult, txError = tx.TExec("delete from user")
	FatalError(t, txError)
	tx.TCommit()
	user_count_commit, uc_error_commit := orm.Fork().Count(u)
	FatalError(t, uc_error_commit)
	if user_count_commit != 0 {
		t.Fatal("can not commit")
	}

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

	/*
		count, err = orm.Fork().Where("`id` < ?", 3).Count(u)
		if err != nil {
			t.Fatal("count where error", err)
		}
		if count != 2 {
			t.Fatal("count where is not one")
		}
		t.Log("count where:", count)
	*/

}

func TestPage(t *testing.T) {
	orm := getOrm()
	u := new(user)

	insert_count := toolkit.RandInt(5, 10)
	t.Log("insert records", insert_count)
	for i := 0; i < insert_count; i++ {
		u := new(user)
		u.Name = toolkit.RandomString(5)
		err := orm.Fork().Save(u)
		FatalError(t, err)
	}

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

func TestWhereIn(t *testing.T) {
	orm := getOrm()
	u := new(user)
	//count, _ := orm.Count(u)
	count := 1000

	uids := make([]interface{}, count)
	for i, _ := range uids {
		uids[i] = count - i
	}

	us_err, us := orm.WhereIn("`id`", uids...).All(u)
	if us_err != nil {
		t.Fatal(us_err)
	}
	t.Log(us)
}

/*
func TestProcedure(t *testing.T) {
	orm := getOrm()
	rows, err := orm.Query("call get_user_by_id(?)", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(rows)

}
*/

func getOrm() *Orm {
	orm, _ := NewOrm("mysql", "127.0.0.1", "root", "123456", "test", 3306)
	return orm
}
