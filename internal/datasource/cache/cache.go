// TODO Вы можете редактировать этот файл по вашему усмотрению

package cache

import (
	"context"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/internal/datasource"
	"time"
)

type Client struct {
	//conn   *bigcache.BigCache
	source datasource.Datasource
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Get(ctx context.Context, key string) (any, error) {
	//TODO implement me
	panic("implement me")
}
