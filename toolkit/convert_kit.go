package toolkit

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
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
		fmt.Println(buf.String())
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

func ZipObj(o interface{}, re map[string]interface{}) (err error) {
	defer func() {
		if zipErr := recover(); zipErr != nil {
			err = errors.New("zip error")
		}
	}()

	val := reflect.ValueOf(o).Elem()
	modelType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		f := modelType.Field(i)
		re[f.Name] = val.Field(i).Interface()
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

func ConvertToString(o interface{}) string {
	switch o.(type) {
	case string:
		return o.(string)
	case int:
		return strconv.Itoa(o.(int))
	case int64:
		return fmt.Sprintf("%d", o.(int64))
	}
	return ""
}

func ConvertToInt(o interface{}) int {
	switch o.(type) {
	case int64:
		return int(o.(int64))
	case int:
		return o.(int)
	case string:
		i, _ := strconv.Atoi(o.(string))
		return i
	}
	return 0
}

func ConvertToInt64(o interface{}) int64 {
	switch o.(type) {
	case int64:
		return o.(int64)
	case int:
		return int64(o.(int))
	case string:
		re, err := strconv.ParseInt(o.(string), 10, 0)
		if err != nil {
			return 0
		} else {
			return re
		}
	}
	return 0
}
