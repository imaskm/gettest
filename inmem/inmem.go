package inmem

import (
	"errors"
	"sync"

	"github.com/imaskm/getir/types"
)

var store sync.Map

func Set(data *types.Data) {
	store.Store(data.Key, data.Value)
}

func Get(key interface{}) (*types.Data, error) {
	value, ok := store.Load(key)
	if ok {
		return &types.Data{Key: key.(string), Value: value.(string)}, nil
	}

	return nil, errors.New("not found")
}
