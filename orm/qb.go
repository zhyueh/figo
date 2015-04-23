package orm

type QueryBuilder interface {
	Select(fields ...string) QueryBuilder
	From(tables ...string) QueryBuilder
	InnerJoin(table string) QueryBuilder
	LeftJoin(table string) QueryBuilder
	RightJoin(table string) QueryBuilder
	On(cond string) QueryBuilder
	Where(cond string) QueryBuilder
	And(cond string) QueryBuilder
	Or(cond string) QueryBuilder
	In(vals ...string) QueryBuilder
	OrderBy(fields ...string) QueryBuilder
	Asc() QueryBuilder
	Desc() QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	GroupBy(fields ...string) QueryBuilder
	Having(cond string) QueryBuilder
	Update(tables ...string) QueryBuilder
	Set(kv ...string) QueryBuilder
	Delete(tables ...string) QueryBuilder
	InsertInto(table string, fields ...string) QueryBuilder
	Values(vals ...string) QueryBuilder
	Subquery(sub string, alias string) string
	String() (string, []interface{})
}
