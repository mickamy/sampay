package kvs

import "errors"

var (
	ErrKeyNotFound                = errors.New("kvs: key not found")
	ErrFailedToMarshalOnMemoize   = errors.New("kvs: marshal on memoize failed")
	ErrFailedToSetOnMemoize       = errors.New("kvs: set on memoize failed")
	ErrFailedToUnmarshalOnMemoize = errors.New("kvs: unmarshal on memoize failed")
)
