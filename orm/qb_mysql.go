package orm

import (
	"fmt"
	"strconv"
	"strings"
)

const MYSQL_COMMA_SPACE = ", "

//mysql qb for mysql quuery builder
type MySQLQB struct {
	fields []string
	values []interface{}

	keyFields []string
	keyValues []interface{}

	tables     []string
	conditions string
	orders     []string
	limit      string
}

func NewMySQLQB() *MySQLQB {
	re := new(MySQLQB)
	re.fields = make([]string, 0)
	re.values = make([]interface{}, 0)
	re.tables = make([]string, 0)
	re.conditions = ""
	re.orders = make([]string, 0)
	re.limit = ""

	return re
}

func (this *MySQLQB) KeyFields(fields []string) {
	this.keyFields = fields
}

func (this *MySQLQB) KeyValues(values []interface{}) {
	this.keyValues = values
}

func (this *MySQLQB) Fields(fields []string) {
	this.fields = fields
}

func (this *MySQLQB) Values(values []interface{}) {
	this.values = values
}

func (this *MySQLQB) Table(table string) {
	this.tables = append(this.tables, table)
}

func (this *MySQLQB) Where(cond string, i ...interface{}) {
	this.conditions = cond
	this.values = append(this.values, i...)
}

func (this *MySQLQB) And(cond string, i interface{}) {
	this.conditions = fmt.Sprintf("%s and %s", this.conditions, cond)
	this.values = append(this.values, i)
}

func (this *MySQLQB) Or(cond string, i interface{}) {
	this.conditions = fmt.Sprintf("%s or %s", this.conditions, cond)
	this.values = append(this.values, i)
}

func (this *MySQLQB) Order(order string) {
	this.orders = append(this.orders, order)
}

func (this *MySQLQB) Limit(limit int) {
	this.limit = strconv.Itoa(limit)
}

func (this *MySQLQB) Page(index, num int) {
	this.limit = fmt.Sprintf("%d, %d", index*num, num)
}

func (this *MySQLQB) Count() (string, []interface{}) {
	sql := fmt.Sprintf(
		"SELECT count(1) as num FROM %s %s",
		this.getTablesString(),
		this.getWhereString(),
	)

	return sql, this.values

}

func (this *MySQLQB) Select() (string, []interface{}) {
	sql := fmt.Sprintf(
		"SELECT %s FROM %s %s %s %s",
		this.getFieldsString(),
		this.getTablesString(),
		this.getWhereString(),
		this.getOrderString(),
		this.getLimitString(),
	)

	return sql, this.values
}

func (this *MySQLQB) Update() (string, []interface{}) {
	sql := fmt.Sprintf(
		"update %s set %s %s",
		this.getTablesString(),
		this.getFieldValuePair(),
		this.getWhereString(),
	)
	return sql, this.values
}

func (this *MySQLQB) Insert() (string, []interface{}) {
	sql := fmt.Sprintf(
		"insert into %s (%s) value(%s)",
		this.getTablesString(),
		this.getFieldsString(),
		this.getValuesSpace(),
	)
	return sql, this.values
}

func (this *MySQLQB) InsertIgnore() (string, []interface{}) {
	if len(this.keyValues) == 0 {
		return this.Insert()
	}

	vals := make([]interface{}, len(this.keyFields))
	for i, v := range this.keyValues {
		vals[i] = v
	}
	vals = append(vals, this.values...)
	vals = append(vals, this.values...)

	sql := fmt.Sprintf(
		"insert into %s (%s, %s) value(%s, %s) on duplicate key update %s",
		this.getTablesString(),
		this.getKeyFieldsString(),
		this.getFieldsString(),
		this.getKeyValuesSpace(),
		this.getValuesSpace(),
		this.getFieldValuePair())

	return sql, vals

}

func (this *MySQLQB) getFieldValuePair() string {
	fv := make([]string, len(this.fields))

	for i, v := range this.fields {
		fv[i] = fmt.Sprintf("`%s`=?", v)
	}
	return strings.Join(fv, MYSQL_COMMA_SPACE)
}

func (this *MySQLQB) getValuesSpace() string {
	s := make([]string, len(this.values))
	for i, _ := range this.values {
		s[i] = "?"
	}

	return strings.Join(s, MYSQL_COMMA_SPACE)
}

func (this *MySQLQB) getKeyValuesSpace() string {
	s := make([]string, len(this.keyValues))
	for i, _ := range this.keyValues {
		s[i] = "?"
	}

	return strings.Join(s, MYSQL_COMMA_SPACE)
}

func (this *MySQLQB) getKeyFieldsString() string {
	if len(this.keyFields) == 0 {
		return "*"
	} else {
		fields := make([]string, len(this.keyFields))
		for i, v := range this.keyFields {
			fields[i] = fmt.Sprintf("`%s`", v)
		}

		return strings.Join(fields, MYSQL_COMMA_SPACE)
	}
}

func (this *MySQLQB) getFieldsString() string {
	if len(this.fields) == 0 {
		return "*"
	} else {
		fields := make([]string, len(this.fields))
		for i, v := range this.fields {
			fields[i] = fmt.Sprintf("`%s`", v)
		}

		return strings.Join(fields, MYSQL_COMMA_SPACE)
	}
}

func (this *MySQLQB) getTablesString() string {
	return strings.Join(this.tables, MYSQL_COMMA_SPACE)
}

func (this *MySQLQB) getWhereString() string {
	if len(this.conditions) == 0 {
		return ""
	} else {
		return fmt.Sprintf(" WHERE %s", this.conditions)
	}
}

func (this *MySQLQB) getOrderString() string {
	order := strings.Join(this.orders, MYSQL_COMMA_SPACE)
	if len(order) > 0 {
		return fmt.Sprintf("order by %s", order)
	}
	return ""
}

func (this *MySQLQB) getLimitString() string {
	if len(this.limit) > 0 {
		return fmt.Sprintf("LIMIT %s", this.limit)
	}
	return ""
}
