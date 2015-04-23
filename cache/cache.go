package cache

import (
	"github.com/xuyu/goredis"
	"strconv"
	"time"
)

const (
	NoKey = 1
)

type Cache struct {
	redisClient *goredis.Redis
}

func NewCache(addr, password, db string) (*Cache, error) {
	dbid, dbiderr := strconv.Atoi(db)
	if dbiderr != nil {
		return nil, dbiderr
	}

	config := &goredis.DialConfig{
		Network:  "tcp",
		Address:  addr,
		Database: dbid,
		Password: password,
		Timeout:  10 * time.Second,
		MaxIdle:  100,
	}
	redisClient, err := goredis.Dial(config)
	if err != nil {
		return nil, err
	}

	return &Cache{redisClient: redisClient}, nil
}

func (this *Cache) Close() {
	this.redisClient.ClosePool()
}

func (this *Cache) Exists(key string) (bool, error) {
	return this.redisClient.Exists(key)
}

func (this *Cache) GetValue(key string) (string, error) {
	data, err := this.redisClient.Get(key)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (this *Cache) SetValue(key, value string) error {
	return this.redisClient.SimpleSet(key, value)
}

func (this *Cache) SetValueEx(key, value string, seconds int) error {
	return this.redisClient.Set(key, value, seconds, 0, false, false)
}

func (this *Cache) SetTTL(key string, seconds int) error {
	_, err := this.redisClient.Expire(key, seconds)
	return err
}
