package cache2go

import (
	"errors"
)

var (
	ErrKeyNotFound           = errors.New("Key not found in cache")
	ErrKeyNotFoundOrLoadable = errors.New("Key not found and could not be loaded into cache")
)
