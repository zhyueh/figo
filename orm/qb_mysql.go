package orm

import (
	"fmt"
	"strings"
)

const MYSQL_COMMA_SPACE = ", "

//mysql qb for mysql quuery builder
type MySQLQB struct {
	fields     []string
	values     []interface{}
	tables     []string
	conditions string
	orders     []string
	limit      int
}

func NewMySQLQB() *MySQLQB {
	re := new(MySQLQB)
	re.fields = make([]string, 0)
	re.values = make([]interface{}, 0)
	re.tables = make([]string, 0)
	re.conditions = ""
	re.orders = make([]string, 0)
	re.limit = 0

	return re
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

func (this *MySQLQB) Where(cond string, i interface{}) {
	this.conditions = cond
	args := make([]interface{}, 1)
	args[0] = i
	this.values = args
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
	this.limit = limit
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
	return strings.Join(this.orders, MYSQL_COMMA_SPACE)
}

func (this *MySQLQB) getLimitString() string {
	if this.limit > 0 {
		return fmt.Sprintf("LIMIT %d", this.limit)
	}
	return ""
}
