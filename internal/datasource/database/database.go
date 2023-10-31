// TODO Вы можете редактировать этот файл по вашему усмотрению

package database

import (
	"context"
	"errors"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/pkg"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	DB pkg.Database
}

func NewDatabaseClient(fileName string) (*Client, error) {
	file, err := Create(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	client := Client{
		DB: pkg.Database{
			FileName: fileName,
		},
	}

	return &client, nil
}

func Create(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
			return nil, err // log.Println(err)
		}
		file, err = os.Create(fileName)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return file, nil
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := c.DB.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = c.DB.Exec(ctx, "INSERT INTO dictionary ( key , value ) VALUES ( $1 , $2 ) ;", key, value)
	if err != nil {
		c.DB.Rollback(ctx)
		return err
	}

	err = c.DB.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string) (any, error) {
	err := c.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}

	dict, err := c.DB.QueryRow(ctx, "SELECT * FROM dictionary WHERE key = $1 ;", key)
	if err != nil {
		c.DB.Rollback(ctx)
		return nil, err
	}

	err = c.DB.Commit(ctx)
	if err != nil {
		return nil, err
	}

	var value string
	err = dict.Scan(&value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
