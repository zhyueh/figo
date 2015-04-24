package toolkit

import (
	"fmt"
	"strconv"
)

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
