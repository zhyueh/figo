package toolkit

import (
	"strconv"
	"testing"
)

type user struct {
	Id   int
	Name string
}

func TestObject(t *testing.T) {
	u := user{}
	u.Id = RandInt(0, 100)
	u.Name = RandomString(100)

	b, err := ObjectToByte(u)
	if err != nil {
		t.Fatal("object encode error", err)
	}

	newu := user{}
	err = ByteToObject(b, &newu)
	if err != nil {
		t.Fatal("object decode error", err)
	}

	if u.Id != newu.Id || u.Name != newu.Name {
		t.Fatal("error object kit")
	}

}

func TestMSI(t *testing.T) {
	msi := make(map[string]interface{}, 10)
	nmsi := make(map[string]interface{}, 10)

	for i := 0; i < 10; i++ {
		seed := RandInt(0, 10)
		if seed%2 == 0 {
			msi[strconv.Itoa(seed)] = seed
		} else {
			msi[strconv.Itoa(seed)] = RandomString(10)
		}
	}

	b, err := ObjectToByte(msi)
	if err != nil {
		t.Fatal(err)
	}

	err = ByteToObject(b, &nmsi)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range msi {
		if val, exists := nmsi[k]; exists &&
			ConvertToString(v) == ConvertToString(val) {
			t.Log("the same", k, v)
		} else {
			t.Fatal("error", k, v)
		}
	}
}
