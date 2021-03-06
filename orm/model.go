package orm

import (
	//"fmt"
	"github.com/zhyueh/figo/toolkit"
	"reflect"
	"time"
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
			if auto == 0 {
				return true
			}
		}
	}

	return false
}

func GetKeyFieldValues(o interface{}) ([]string, []interface{}, bool) {
	val := reflect.ValueOf(o).Elem()
	modelType := val.Type()

	fields := make([]string, 0)
	vals := make([]interface{}, 0)

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		if f.Tag.Get("orm") == "auto" || f.Tag.Get("orm") == "primary" {
			//return ModelFieldToSqlField(f), val.Field(i).Interface(), true
			fields = append(fields, ModelFieldToSqlField(f))
			vals = append(vals, val.Field(i).Interface())
		}
	}
	return fields, vals, true
	//return "", nil, false
}

func GetIdFieldValue(o interface{}) (string, interface{}, bool) {
	val := reflect.ValueOf(o).Elem()
	modelType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		if f.Tag.Get("orm") == "auto" || f.Tag.Get("orm") == "primary" {
			return ModelFieldToSqlField(f), val.Field(i).Interface(), true
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
		if ormTag != "" && ormTag != "auto" && ormTag != "primary" {

			if ignoreField(f, val.Field(i)) {
				continue
			}

			fields = append(fields, ModelFieldToSqlField(f))
			values = append(values, updateFieldValue(
				f,
				val.Field(i),
			))
		}
	}

	return fields, values

}

func updateFieldValue(f reflect.StructField, v reflect.Value) interface{} {
	if f.Type.PkgPath() == "time" && f.Type.Name() == "Time" {
		//return different  value when orm is in ("int", "date")
		//TODO is a pointer?
		t := v.Interface().(time.Time)
		if f.Tag.Get("orm") == "int" {
			return t.Unix()
		} else if f.Tag.Get("orm") == "date" {
			return t.Format("2006-01-02")
		}
	}

	return v.Interface()

}

func ignoreField(field reflect.StructField, fieldVal reflect.Value) bool {
	if fieldVal.Kind() == reflect.String {
		if fieldVal.String() == "" && field.Tag.Get("empty") == "ignore" {
			return true
		}
	}

	if field.Tag.Get("readonly") != "" {
		return true
	}
	return false
}

func ModelTableName(model ModelInterface) string {
	if name := model.TableName(); name != DefaultTable {
		return name
	} else {
		modelType := reflect.ValueOf(model).Elem().Type()
		modelName := modelType.Name()
		return toolkit.CamelCaseToUnderScore(modelName)
	}
}

func ModelFieldToSqlField(field reflect.StructField) string {
	if name := field.Tag.Get("name"); name != "" {
		return name
	}

	return toolkit.CamelCaseToUnderScore(field.Name)
}

func ModelUpdateId(model interface{}, id int64) {
	val := reflect.ValueOf(model).Elem()
	modelType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		if f.Tag.Get("orm") == "auto" &&
			(f.Type.Kind() == reflect.Int ||
				f.Type.Kind() == reflect.Int64) {
			val.FieldByName(f.Name).SetInt(id)
			return
		}
	}

}

func DbRowToModelEx(row DbRow, model interface{}) {
	//fmt.Println(row)
	val := reflect.ValueOf(model).Elem()
	modelType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		if f.Tag.Get("orm") == "extend" && f.Type.Kind() == reflect.Struct {
			DbRowToModelEx(row, val.Field(i).Addr().Interface())
		} else if name, canSet := getCanSetFieldName(f); canSet {
			o, exists := row[name]
			if exists {
				val.Field(i).Set(convertSqlValueToFieldValue(o, f))
			} else {
				//fmt.Println("do not exists ", name)
			}
		} else {
			//fmt.Println("can not set")
		}
	}
}

func DbRowToModel(row DbRow, model interface{}) {
	//fmt.Println(row)
	val := reflect.ValueOf(model).Elem()
	modelType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		if name, canSet := getCanSetFieldName(f); canSet {
			o, exists := row[name]
			if exists {
				val.Field(i).Set(convertSqlValueToFieldValue(o, f))
			} else {
				//fmt.Println("do not exists ", name)
			}
		} else {
			//fmt.Println("can not set")
		}
	}
}

func convertSqlValueToFieldValue(o interface{}, f reflect.StructField) reflect.Value {
	if f.Type.Kind() == reflect.String {
		return reflect.ValueOf(toolkit.ConvertToString(o))
	} else if f.Type.Kind() == reflect.Int {
		return reflect.ValueOf(toolkit.ConvertToInt(o))
	} else if f.Type.Kind() == reflect.Float64 {
		return reflect.ValueOf(toolkit.ConvertToFloat64(o))
	} else if f.Type.Kind() == reflect.Struct {
		//convert time
		if f.Type.PkgPath() == "time" && f.Type.Name() == "Time" {
			return reflect.ValueOf(toolkit.ConvertToTime(o))

		}
	}
	return reflect.Zero(f.Type)
}

func getCanSetFieldName(f reflect.StructField) (string, bool) {
	ormTag := f.Tag.Get("orm")
	if ormTag != "" {
		return ModelFieldToSqlField(f), true
	}
	return "", false
}
