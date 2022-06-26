package mem

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/twistedogic/store"
)

type Store struct {
	m *sync.Map
}

func New() (Store, error) {
	return Store{m: new(sync.Map)}, nil
}

func (s Store) Get(ctx context.Context, key []byte) (store.Item, bool) {
	v, ok := s.m.Load(string(key))
	if !ok {
		return store.Item{}, ok
	}
	i, ok := v.(store.Item)
	return i, ok
}

func (s Store) Set(ctx context.Context, i store.Item) error {
	s.m.Store(string(i.Key), i)
	return nil
}

func (s Store) Delete(ctx context.Context, key []byte) error {
	if _, exist := s.Get(ctx, key); !exist {
		return fmt.Errorf("key %s not found", key)
	}
	s.m.Delete(string(key))
	return nil
}

func (s Store) PrefixScan(ctx context.Context, prefix []byte) ([]store.Item, error) {
	out := make([]store.Item, 0)
	s.m.Range(func(key, val interface{}) bool {
		if strings.HasPrefix(key.(string), string(prefix)) {
			out = append(out, val.(store.Item))
		}
		return true
	})
	return out, nil
}
