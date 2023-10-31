package pkg

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/pkg/sql_parser"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Database struct {
	FileName       string
	backUpFileName string
	sync.Mutex
}

func (d *Database) Begin(ctx context.Context) error {
	d.Lock()

	sourceFile, err := os.Open(d.FileName)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	fileNameWithoutExt := strings.TrimSuffix(d.FileName, filepath.Ext(d.FileName))
	sliceToCreateTempFileName := []string{fileNameWithoutExt, "_temp", filepath.Ext(d.FileName)}
	backUpFileName := strings.Join(sliceToCreateTempFileName, "")
	backUpFile, err := os.Create(backUpFileName)
	if err != nil {
		return err
	}
	defer backUpFile.Close()

	d.backUpFileName = backUpFileName

	_, err = io.Copy(backUpFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) Commit(ctx context.Context) error {
	err := os.Remove(d.backUpFileName)
	if err != nil {
		return err
	}
	defer d.Unlock()
	return nil
}

func (d *Database) Rollback(ctx context.Context) error {
	err := os.Remove(d.FileName)
	if err != nil {
		return err
	}

	err = os.Rename(d.backUpFileName, d.FileName)
	if err != nil {
		return err
	}
	d.backUpFileName = ""
	defer d.Unlock()
	return nil
}

func (d *Database) Exec(ctx context.Context, sql string, args ...any) (commandTag pgconn.CommandTag, err error) {
	queryInfo, err := sql_parser.Parse(sql)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	queryRequirement := sql_parser.QueryInfo{
		Operation:     "INSERT",
		Table:         "dictionary",
		IdentSlice:    []string{"key", "value"},
		TemplateSlice: []string{"$1", "$2"},
	}

	if queryInfo.Equals(queryRequirement) {
		if len(args) < 2 {
			return pgconn.CommandTag{}, ArgsNotSpecifiedError
		}
		err := d.Insert(args[0].(string), args[1])
		if err != nil {
			return pgconn.CommandTag{}, err
		}
		return pgconn.NewCommandTag("INSERT"), nil
	} else {
		return pgconn.CommandTag{}, SqlScriptNotSupportedError
	}
}

func (d *Database) QueryRow(ctx context.Context, sql string, args ...any) (pgx.Row, error) {
	queryInfo, err := sql_parser.Parse(sql)
	if err != nil {
		return nil, err
	}

	queryRequirement := sql_parser.QueryInfo{
		Operation:     "SELECT",
		Table:         "dictionary",
		IdentSlice:    []string{"key"},
		TemplateSlice: []string{"$1"},
	}

	if queryInfo.Equals(queryRequirement) {
		if len(args) < 1 {
			return nil, ArgsNotSpecifiedError
		}
		key := args[0].(string)
		value, err := d.Get(key)
		if err != nil {
			return nil, err
		}
		return &Dict{
			Key:   key,
			Value: value,
		}, nil
	} else {
		return nil, SqlScriptNotSupportedError
	}
}

// Insert and checks if key is not repeatable
func (d *Database) Insert(key string, value any) error {
	listDict, err := d.ListDict()
	if err != nil {
		return err
	}

	if listDict.IsExist(key) {
		return KeyAlreadyExist
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
	return nil, KeyNotFound

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
