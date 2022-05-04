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

func (s Store) Get(ctx context.Context, key string) (store.Item, error) {
	v, ok := s.m.Load(key)
	if !ok {
		return store.Item{}, fmt.Errorf("key %s not found", key)
	}
	i, ok := v.(store.Item)
	if !ok {
		return store.Item{}, fmt.Errorf("value for key %s is not Item", key)
	}
	return i, nil
}

func (s Store) Exists(ctx context.Context, key string) (bool, error) {
	_, ok := s.m.Load(key)
	return ok, nil
}

func (s Store) Set(ctx context.Context, i store.Item) error {
	s.m.Store(i.Key, i)
	return nil
}

func (s Store) Delete(ctx context.Context, key string) error {
	exist, err := s.Exists(ctx, key)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("key %v not found", key)
	}
	s.m.Delete(key)
	return nil
}

func (s Store) PrefixScan(ctx context.Context, prefix string) ([]store.Item, error) {
	out := make([]store.Item, 0)
	s.m.Range(func(key, val interface{}) bool {
		if strings.HasPrefix(key.(string), prefix) {
			out = append(out, val.(store.Item))
		}
		return true
	})
	return out, nil
}
