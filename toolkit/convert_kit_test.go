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
