package common

import (
	"sync"
)

// 全局数据缓存
type Cache struct {
	sync.RWMutex
	Values map[string]interface{}
}

var cache *Cache

func GetCache() *Cache {
	if cache == nil {
		cache = &Cache{
			Values: make(map[string]interface{}),
		}
	}
	return cache
}

// 设置缓存
func (c *Cache) Set(k string, v interface{}) {
	c.Lock()
	defer c.Unlock()
	c.Values[k] = v
}

// 获取缓存
func (c *Cache) Get(k string) interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.Values[k]
}

// 获取int
func (c *Cache) GetInt(k string) int {
	c.RLock()
	defer c.RUnlock()

	v := c.Values[k]
	if v == nil {
		return 0
	}
	res, ok := v.(int)
	if !ok {
		return 0
	}
	return res
}

// 获取int64
func (c *Cache) GetInt64(k string) int64 {
	c.RLock()
	defer c.RUnlock()

	v := c.Values[k]
	if v == nil {
		return 0
	}
	res, ok := v.(int64)
	if !ok {
		return 0
	}
	return res
}

// 获取uint64
func (c *Cache) GetUInt64(k string) uint64 {
	c.RLock()
	defer c.RUnlock()

	v := c.Values[k]
	if v == nil {
		return 0
	}
	res, ok := v.(uint64)
	if !ok {
		return 0
	}
	return res
}

// 获取int
func (c *Cache) GetFloat(k string) float64 {
	c.RLock()
	defer c.RUnlock()

	v := c.Values[k]
	if v == nil {
		return 0
	}
	res, ok := v.(float64)
	if !ok {
		return 0
	}
	return res
}

// 获取bool
func (c *Cache) GetBool(k string) bool {
	c.RLock()
	defer c.RUnlock()

	v := c.Values[k]
	if v == nil {
		return false
	}
	res, ok := v.(bool)
	if !ok {
		return false
	}
	return res
}

// 获取string
func (c *Cache) GetString(k string) string {
	c.RLock()
	defer c.RUnlock()

	v := c.Values[k]
	if v == nil {
		return ""
	}
	res, ok := v.(string)
	if !ok {
		return ""
	}
	return res
}
