package orm

import (
	"fmt"
	"github.com/zhyueh/figo/toolkit"
	"reflect"
)

type ModelInterface interface {
	TableName() string
}

const DefaultTable = "==use=default=table=name=="

type ModelSample struct {
	id   int    `name:"id" orm:"auto"`
	name string `name:"name" orm:"varchar"`
}

type Model struct {
}

func (this *Model) TableName() string {
	return DefaultTable
}

func NeedInsertModel(model interface{}) bool {
	val := reflect.ValueOf(model).Elem()
	modelType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		if f.Tag.Get("orm") == "auto" {
			auto := val.Field(i).Int()
			if auto > 0 {
				return false
			}
		}
	}

	return true
}

func GetIdFieldValue(o interface{}) (string, interface{}, bool) {
	val := reflect.ValueOf(o).Elem()
	modelType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		if f.Tag.Get("orm") == "auto" {
			return f.Tag.Get("name"), val.Field(i).Interface(), true
		}
	}

	return "", nil, false
}

func GetSaveModelFieldValues(model ModelInterface) ([]string, []interface{}) {

	val := reflect.ValueOf(model).Elem()
	modelType := val.Type()

	fields := make([]string, 0)
	values := make([]interface{}, 0)

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		ormTag := f.Tag.Get("orm")
		if ormTag != "" && ormTag != "auto" {

			if ignoreField(f, val.Field(i)) {
				continue
			}

			fields = append(fields, ModelFieldToSqlField(f))
			values = append(values, val.Field(i).Interface())
		}
	}

	return fields, values

}

func ignoreField(field reflect.StructField, fieldVal reflect.Value) bool {
	if fieldVal.Kind() == reflect.String {
		if fieldVal.String() == "" && field.Tag.Get("empty") == "ignore" {
			return true
		}
	}
	return false
}

func ModelTableName(model ModelInterface) string {
	if name := model.TableName(); name != DefaultTable {
		return name
	} else {
		modelType := reflect.ValueOf(model).Elem().Type()
		modelName := modelType.Name()
		return modelName
	}
}

func ModelFieldToSqlField(field reflect.StructField) string {
	if name := field.Tag.Get("name"); name != "" {
		return name
	}

	return ""
}

func ModelUpdateId(model interface{}, id int64) {
	val := reflect.ValueOf(model).Elem()
	modelType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		if f.Tag.Get("orm") == "auto" {
			return
		}
	}

}

func DbRowToModel(row DbRow, model interface{}) {
	fmt.Println(row)
	val := reflect.ValueOf(model).Elem()
	modelType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		if name, canSet := getCanSetFieldName(f); canSet {
			o, exists := row[name]
			if exists {
				val.Field(i).Set(convertSqlValueToFieldValue(o, f))
			} else {
				fmt.Println("do not exists ", name)
			}
		} else {
			fmt.Println("can not set")
		}
	}
}

func convertSqlValueToFieldValue(o interface{}, f reflect.StructField) reflect.Value {
	if f.Type.Kind() == reflect.String {
		return reflect.ValueOf(toolkit.ConvertToString(o))
	} else if f.Type.Kind() == reflect.Int {
		return reflect.ValueOf(toolkit.ConvertToInt(o))
	}
	return reflect.ValueOf("")
}

func getCanSetFieldName(f reflect.StructField) (string, bool) {
	ormTag := f.Tag.Get("orm")
	if ormTag != "" {
		return f.Tag.Get("name"), true
	}
	return "", false
}
