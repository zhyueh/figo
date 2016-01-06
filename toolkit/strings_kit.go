package toolkit

import (
	"strings"
)

func SplitString(s, sep string) []string {
	re := make([]string, 0)
	for _, v := range strings.Split(s, sep) {
		if len(v) != 0 {
			re = append(re, v)
		}
	}

	return re
}
