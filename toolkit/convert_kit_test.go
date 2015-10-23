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

func TestToFloat64(t *testing.T) {
	var i = 1.0
	in := []interface{}{1.0, 1, "1", "1.000"}
	for _, v := range in {
		if ConvertToFloat64(v) != i {
			t.Fatal("failed to float64", v)
		}
	}
}

func TestToString(t *testing.T) {
	var i = "123"
	var ui uint8
	ui = 123
	in := []interface{}{123, ui, "123"}
	for _, v := range in {
		if ConvertToString(v) != i {
			t.Fatal("failed to string", v)
		}
	}

	if tmp := ConvertToString(123.4); tmp != "123.4" {
		t.Fatal("failed to string", 123.4, tmp)
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
