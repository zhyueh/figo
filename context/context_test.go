package context

import (
	"net/http"
	"testing"
)

func TestFunc(t *testing.T) {

	assertEqual := func(v1, v2 interface{}) {
		if v1 != v2 {
			t.Errorf("expect %v but %v", v1, v2)
		}
	}

	r1, _ := http.NewRequest("GET", "http://127.0.0.1", nil)
	r2, _ := http.NewRequest("GET", "http://127.0.0.1:8080", nil)
	r3, _ := http.NewRequest("POST", "http://127.0.0.1:9090", nil)

	Set(r1, "key", "1")
	Set(r2, "key", "2")
	assertEqual(Get(r2, "key"), "2")
	assertEqual(GetString(r2, "key"), "2")
	assertEqual(Get(r1, "key"), "1")

	assertEqual(len(GetAll(r1)), 1)
	Delete(r1, "key")
	assertEqual(len(GetAll(r1)), 0)

	assertEqual(len(GetAll(r2)), 1)
	assertEqual(RequestNum(), 2)
	Clear(r2)
	assertEqual(RequestNum(), 1)
	assertEqual(Get(r3, "key"), nil)

}
