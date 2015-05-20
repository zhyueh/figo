package toolkit

import (
	"testing"
)

func TestH(t *testing.T) {

	paras := make(map[string]string, 0)
	paras["bbb"] = "BBB"
	paras["ccc"] = "111"
	template := "aaa{{bbb|Lower}}123{ {{ccc}}"
	tmp := NewTemplate(template, paras)

	result := tmp.Do()
	if result != "aaabbb123{ 111" {
		t.Fatal("template not work", result)
	}

}
