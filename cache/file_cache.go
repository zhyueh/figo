package cache

import (
	"github.com/zhyueh/figo/toolkit"
)

type FileCache struct {
	path string
}

func NewFileCache(path string) *FileCache {
	re := new(FileCache)
	re.path = path
	return re
}

func PathOfKey(category, key string) string {
	md5 := toolkit.MakeMd5Str([]byte(key))
	return toolkit.PathJoin(
		category,
		md5[0:2],
		md5,
	)
}

func (this *FileCache) FullPath(category, key string) string {
	return toolkit.PathJoin(
		this.path,
		PathOfKey(category, key),
	)
}

func (this *FileCache) Exists(category, key string) bool {
	return toolkit.IsExists(this.FullPath(category, key))
}

func (this *FileCache) Get(category, key string) ([]byte, error, bool) {
	if exists := this.Exists(category, key); !exists {
		return nil, nil, false
	}
	data, err := toolkit.ReadAll(this.FullPath(category, key))
	return data, err, true
}

func (this *FileCache) Set(category, key string, data []byte) error {
	fullpath := this.FullPath(category, key)
	toolkit.EnsureDirExists(toolkit.ParentPath(fullpath))

	return toolkit.WriteAll(data, fullpath)
}

func (this *FileCache) Del(category, key string) error {
	toolkit.RemoveFileIfExists(this.FullPath(category, key))
	return nil
}

func (this *FileCache) GetEx(category, key string, i interface{}) (error, bool) {
	data, err, exists := this.Get(category, key)
	if err == nil && exists {
		err = toolkit.ByteToObject(data, i)
		return err, true
	} else {
		return err, exists
	}
}

func (this *FileCache) SetEx(category, key string, i interface{}) error {
	data, err := toolkit.ObjectToByte(i)
	if err != nil {
		return err
	}

	return this.Set(category, key, data)

}
