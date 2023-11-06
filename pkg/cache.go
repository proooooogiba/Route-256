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
	mutex sync.RWMutex
}

type CacheItem struct {
	value      any
	expiration time.Time
}

func (c *Cache) Get(key string) *redis.StringCmd {
	c.update()
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if cacheItem, ok := c.Items[key]; ok {
		if cacheItem.expiration.After(time.Now()) {
			strCmd := redis.NewStringCmd(context.Background())
			strCmd.SetVal(cacheItem.value.(string))
			return strCmd
		}
	}

	strCmd := redis.NewStringCmd(context.Background())
	strCmd.SetErr(CacheNotFoundError)
	return strCmd
}

func (c *Cache) Set(key string, value any, expiration time.Duration) *redis.StatusCmd {
	c.update()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Items[key] = &CacheItem{}
	c.Items[key].value = value
	c.Items[key].expiration = time.Now().Add(expiration)
	return redis.NewStatusCmd(context.Background())
}

func (c *Cache) update() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, cacheItem := range c.Items {
		if cacheItem.expiration.Before(time.Now()) {
			delete(c.Items, key)
		}
	}
}
