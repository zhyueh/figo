package toolkit

import (
	"testing"
)

func TestH(t *testing.T) {

	paras := make(map[string]string, 0)
	paras["bbb"] = "BBB"
	template := "aaa{{bbb|Lower}}"
	tmp := NewTemplate(template, paras)

	result := tmp.Do()
	if result != "aaabbb" {
		t.Fatal("template not work", result)
	}

}
