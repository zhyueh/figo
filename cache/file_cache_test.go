package cache

import (
	"github.com/zhyueh/figo/toolkit"
	"testing"
)

type test_cache struct {
	Id   int
	Name string
}

func TestBaseFileCache(t *testing.T) {
	cache := getFileCache()
	category := "app/install"
	key := toolkit.RandomString(5)
	data := []byte("123")
	cache.Set(category, key, data)
	cacheData, err, exists := cache.Get(category, key)
	if err != nil || !exists {
		t.Fatal("error cache", err, exists)
	}

	if string(cacheData) != string(data) {
		t.Fatal("error data")
	}
}

func TestDelFileCache(t *testing.T) {
	cache := getFileCache()
	category := "app/info"
	key := toolkit.RandomString(10)
	data := []byte("123")
	cache.Set(category, key, data)
	_, _, exists := cache.Get(category, key)
	if !exists {
		t.Fatal("error get/set")
	}
	cache.Del(category, key)
	_, _, exists = cache.Get(category, key)
	if exists {
		t.Fatal("error del")
	}

}

func TestEx(t *testing.T) {
	cache := getFileCache()
	category := "app/info"
	key := toolkit.RandomString(10)
	defer cache.Del(category, key)

	t1 := test_cache{}
	t1.Id = toolkit.RandInt(0, 10)
	t1.Name = toolkit.RandomString(10)
	err := cache.SetEx(category, key, t1)
	if err != nil {
		t.Fatal(err)
	}
	t2 := test_cache{}
	err, _ = cache.GetEx(category, key, &t2)
	if err != nil {
		t.Fatal(err)
	}

	if t1.Id != t2.Id || t2.Name != t2.Name {
		t.Fatal("error ex")
	}

}

func getFileCache() *FileCache {
	toolkit.EnsureDirExists("/tmp/figo_cache")
	return NewFileCache("/tmp/figo_cache")
}
