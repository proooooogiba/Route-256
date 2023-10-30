package pkg

import "errors"

var (
	CacheNotFoundError         = errors.New("key not found in cache")
	NotPointerError            = errors.New("scan target not a pointer")
	SqlScriptNotSupportedError = errors.New("sql script isn't supported")
	ArgsNotSpecifiedError      = errors.New("argument not specified")
	KeyAlreadyExist            = errors.New("dict with same key already exist")
)
