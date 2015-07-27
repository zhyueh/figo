package cache

import (
	"testing"
)

func TestSetEx(t *testing.T) {
	cache, cache_err := getTestCache()
	if cache_err != nil {
		t.Fatal("%v", cache_err)
	}
	for i := 0; i < 10; i++ {
		err := cache.SetValueEx("123", "123", 1000)
		if err != nil {
			t.Fatal("%v", err)
		}
	}

}
func getTestCache() (*Cache, error) {
	return NewCache("127.0.0.1:6379", "", "1")
}
