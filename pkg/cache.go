// TODO Вы можете редактировать этот файл по вашему усмотрению

package pkg

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type Cache struct {
	Items map[string]*CacheItem
	mutex sync.Mutex
}

type CacheItem struct {
	value      any
	expiration time.Time
}

func (c *Cache) Get(key string) *redis.StringCmd {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Update()
	for keyItem, cacheItem := range c.Items {
		if key == keyItem {
			if cacheItem.expiration.After(time.Now()) {
				strCmd := redis.NewStringCmd(context.Background())
				strCmd.SetVal(cacheItem.value.(string))
				return strCmd
			} else {
				delete(c.Items, key)
			}
		}
	}

	strCmd := redis.NewStringCmd(context.Background())
	strCmd.SetErr(CacheNotFoundError)
	return strCmd
}

func (c *Cache) Set(key string, value any, expiration time.Duration) *redis.StatusCmd {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Update()
	c.Items[key] = &CacheItem{}
	c.Items[key].value = value
	c.Items[key].expiration = time.Now().Add(expiration)
	return redis.NewStatusCmd(context.Background(), "SET", key, value)
}

func (c *Cache) Update() {
	for key, cacheItem := range c.Items {
		if cacheItem.expiration.Before(time.Now()) {
			delete(c.Items, key)
		}
	}
}
