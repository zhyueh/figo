package toolkit

import (
	"testing"
)

func TestNamingConvension(t *testing.T) {
	sample := make(map[string]string, 0)
	sample["AaBbCc"] = "aa_bb_cc"

	for k, v := range sample {
		if CamelCaseToUnderScore(k) != v {
			t.Fatal("can not convert", k, "to", v, CamelCaseToUnderScore(k))
		}

		if UnderScoreToCamelCase(v) != k {
			t.Fatal("can not convert", v, "to", k, UnderScoreToCamelCase(v))
		}

	}

	if CamelCaseToUnderScore("AaBbCc") != "aa_bb_cc" {
		t.Fatal("failed AaBbCc")
	}
}

func TestToInt(t *testing.T) {
	var i = 1
	in := []interface{}{1.2, 1, "1", "1.000"}
	for _, v := range in {
		if ConvertToInt(v) != i {
			t.Fatal("failed to int", v)
		}
	}
}

func TestToInt64(t *testing.T) {
	var i int64 = 1
	in := []interface{}{1.2, 1, "1"}
	for _, v := range in {
		if ConvertToInt64(v) != i {
			t.Fatal("failed to int64", v)
		}
	}
}
