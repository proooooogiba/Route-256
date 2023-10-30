// TODO Вы можете редактировать этот файл по вашему усмотрению

package pkg

import (
	"container/list"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type Cache struct {
	items map[string]*list.Element
	list  *list.List
	mutex sync.Mutex
}

type CacheItem struct {
	key        string
	value      interface{}
	expiration time.Time
}

func (c *Cache) Get(key string) *redis.StringCmd {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		item := elem.Value.(*CacheItem)
		if item.expiration.After(time.Now()) {
			c.list.MoveToFront(elem)
			return item.value.(*redis.StringCmd)
		} else {
			c.list.Remove(elem)
			delete(c.items, key)
		}
	}
	return nil
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*CacheItem).value = value
		elem.Value.(*CacheItem).expiration = time.Now().Add(expiration)
		return &redis.StatusCmd{}
	}

	elem := c.list.PushFront(&CacheItem{
		key:        key,
		value:      value,
		expiration: time.Now().Add(expiration),
	})
	c.items[key] = elem
	return &redis.StatusCmd{}
}

//func main() {
//	cache := NewCache()
//
//	// set an item with expiration of 5 seconds
//	cache.Set("key1", "value1", 5*time.Second)
//
//	// get the item before expiration
//	if value, ok := cache.Get("key1"); ok {
//		fmt.Println(value)
//	}
//
//	// wait for the item to expire
//	time.Sleep(6 * time.Second)
//
//	// get the item after expiration
//	if _, ok := cache.Get("key1"); !ok {
//		fmt.Println("item expired")
//	}
//}
