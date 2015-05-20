package toolkit

import (
	"testing"
)

func TestH(t *testing.T) {

	paras := make(map[string]string, 0)
	paras["bbb"] = "BBB"
	template := "aaa{bbb|Lower|Timestamp}"
	tmp := NewTemplate(template, paras)

	if tmp.Do() != "aaa123" {
		t.Fatal("template not work")
	}

}
