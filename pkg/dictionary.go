package pkg

type DictList []*Dict

type Dict struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (d *DictList) IsExist(key string) bool {
	for _, dict := range *d {
		if dict.Key == key {
			return true
		}
	}
	return false
}
