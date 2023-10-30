package pkg

import (
	"reflect"
)

type DictList []*Dict

type Dict struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (d Dict) Scan(dest ...any) error {
	val := reflect.ValueOf(dest[0])
	if val.Kind() != reflect.Ptr {
		return NotPointerError
	}
	val.Elem().Set(reflect.ValueOf(d.Value))
	return nil
}

func (d *DictList) IsExist(key string) bool {
	for _, dict := range *d {
		if dict.Key == key {
			return true
		}
	}
	return false
}
