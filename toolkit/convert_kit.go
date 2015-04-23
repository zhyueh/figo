package toolkit

import (
	"strconv"
)

func ConvertToString(o interface{}) string {
	switch o.(type) {
	case string:
		return o.(string)
	case int:
		return strconv.Itoa(o.(int))
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
	return 0
}
