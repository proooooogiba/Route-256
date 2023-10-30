// TODO Вы можете редактировать этот файл по вашему усмотрению

package cache

import (
	"context"
	"errors"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/internal/datasource"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/internal/datasource/database"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/pkg"
	"time"
)

type Client struct {
	cache pkg.Cache
	db    datasource.Datasource
}

func NewDatabaseWithCacheClient(fileName string) (*Client, error) {
	db, err := database.NewDatabaseClient(fileName)
	if err != nil {
		return nil, err
	}

	return &Client{
		pkg.Cache{
			Items: make(map[string]*pkg.CacheItem),
		},
		db,
	}, nil
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := c.db.Set(ctx, key, value, expiration)
	if err != nil {
		return err
	}

	result := c.cache.Set(key, value, expiration)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (c *Client) Get(ctx context.Context, key string) (any, error) {
	result := c.cache.Get(key)
	if result.Err() != nil {
		if errors.Is(result.Err(), pkg.CacheNotFoundError) {
			value, err := c.db.Get(ctx, key)
			if err != nil {
				return nil, err
			}
			return value, nil
		} else {
			return nil, nil
		}
	}
	value := result.Val()
	return value, nil
}
