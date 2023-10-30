// TODO Вы можете редактировать этот файл по вашему усмотрению

package pkg

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// TODO реализовать только нужные

type Database struct {
	FileName string
}

//func (d *Database) Begin(ctx context.Context) (pgx.Tx, error) {
//
//}
//
//func (d *Database) Commit(ctx context.Context) error {
//
//}
//
//func (d *Database) Rollback(ctx context.Context) error {
//
//}
//
//func (d *Database) Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error) {
//
//}
//
//func (d *Database) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
//
//}
//
//func (d *Database) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
//
//}

// Insert and checks if key is not repeatable
func (d *Database) Insert(key string, value any) error {
	listDict, err := d.ListDict()
	if err != nil {
		return err
	}

	if listDict.IsExist(key) {
		return errors.New("dict with same key already exist")
	}

	newDict := &Dict{
		Key:   key,
		Value: value,
	}

	listDict = append(listDict, newDict)

	bytes, err := json.Marshal(listDict)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(d.FileName, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if _err := file.Truncate(0); _err != nil {
		return _err
	}

	if _, _err := file.Seek(0, 0); _err != nil {
		return _err
	}

	if _, err = file.Write(bytes); err != nil {
		return err
	}

	return nil
}

func (d *Database) Get(key string) (any, error) {
	listDict, err := d.ListDict()
	if err != nil {
		return nil, err
	}

	for _, dict := range listDict {
		if key == dict.Key {
			return dict.Value, nil
		}
	}
	return nil, errors.New("key not found")
}

func (d *Database) ListDict() (DictList, error) {
	file, err := os.OpenFile(d.FileName, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	readAll, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	var list DictList

	if len(readAll) == 0 {
		return list, nil
	}

	if err := json.Unmarshal(readAll, &list); err != nil {
		return nil, err
	}

	return list, nil
}
