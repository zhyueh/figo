package toolkit

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func CamelCaseToUnderScore(code string) string {
	buf := bytes.NewBuffer([]byte{})
	buf.Reset()

	for i, r := range code {

		if r >= 'A' && r <= 'Z' {
			if i != 0 {
				buf.WriteRune('_')
			}
			buf.WriteRune(r + 32)

		} else {
			buf.WriteRune(r)
		}
	}

	return buf.String()

}

func UnderScoreToCamelCase(code string) string {
	buf := bytes.NewBufferString("")

	nextUpper := true
	for _, r := range code {
		if nextUpper {
			buf.WriteRune(ToUpper(r))
			nextUpper = false
		} else if r == '_' {
			nextUpper = true
			continue
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()

}

func ToLower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + 32
	}
	return r
}

func ToUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 32
	}
	return r
}

/*
MSS for map[string]string
MSI for map[string]interface{}

*/

func ZipObjString(o interface{}, re map[string]string) (err error) {
	defer func() {
		if zipErr := recover(); zipErr != nil {
			err = errors.New("zip error")
		}
	}()

	val := reflect.ValueOf(o).Elem()
	modelType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		re[f.Name] = ConvertToString(val.Field(i).Interface())
	}

	err = nil
	return
}

func MSSToMSI(target map[string]string, source map[string]interface{}) {
	for k, v := range source {
		target[k] = ConvertToString(v)
	}
}

func MergeMSS(target, source map[string]string) {
	for k, v := range source {
		target[k] = v
	}
}

func ConvertToTime(o interface{}) time.Time {
	switch o.(type) {
	case string:
		val := o.(string)
		if t, err := time.Parse("2006-01-02", val); err == nil {
			return t
		} else if t, err := time.Parse("2006-01-02 15:04:05", val); err == nil {
			return t
		}
	case int64:
		t := time.Unix(o.(int64), 0)
		return t
	case int:
		t := time.Unix(int64(o.(int)), 0)
		return t

	}
	return time.Time{}

}

func ConvertToFloat64(o interface{}) float64 {
	val := reflect.ValueOf(o)
	kind := reflect.TypeOf(o).Kind()
	switch {
	case kind == reflect.Float32 || kind == reflect.Float64:
		return val.Float()
	case kind == reflect.String:
		if f, err := strconv.ParseFloat(val.String(), 64); err == nil {
			return f
		}
	case kind >= reflect.Int && kind <= reflect.Int64:
		return float64(val.Int())
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return float64(val.Uint())
	}
	return 0.0
}

func ConvertToString(o interface{}) string {
	val := reflect.ValueOf(o)
	kind := reflect.TypeOf(o).Kind()
	switch {
	case kind == reflect.String:
		return val.String()
	case kind == reflect.Float32 || kind == reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64)
	case kind >= reflect.Int && kind <= reflect.Int64:
		return fmt.Sprintf("%d", val.Int())
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return fmt.Sprintf("%d", val.Uint())
	}
	return ""
}

func ConvertToInt(o interface{}) int {
	val := reflect.ValueOf(o)
	kind := reflect.TypeOf(o).Kind()
	switch {
	case kind == reflect.String:
		i, err := strconv.Atoi(val.String())
		if err == nil {
			return i
		}
		f, ferr := strconv.ParseFloat(val.String(), 64)
		if ferr == nil {
			return int(f)
		}
	case kind == reflect.Float32 || kind == reflect.Float64:
		return int(val.Float())
	case kind >= reflect.Int && kind <= reflect.Int64:
		return int(val.Int())
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return int(val.Uint())
	}
	return 0
}

func ConvertToInt64(o interface{}) int64 {
	val := reflect.ValueOf(o)
	kind := reflect.TypeOf(o).Kind()
	switch {
	case kind == reflect.String:
		re, err := strconv.ParseInt(val.String(), 10, 0)
		if err == nil {
			return re
		}

		f, ferr := strconv.ParseFloat(val.String(), 64)
		if ferr == nil {
			return int64(f)
		}
	case kind == reflect.Float32 || kind == reflect.Float64:
		return int64(val.Float())
	case kind >= reflect.Int && kind <= reflect.Int64:
		return val.Int()
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return int64(val.Uint())
	}
	return 0
}
