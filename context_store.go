package store

import "context"

type StoreWithContext struct {
	ctx   context.Context
	store Store
}

func NewStoreWithContext(ctx context.Context, s Store) StoreWithContext {
	return StoreWithContext{ctx: ctx, store: s}
}

func (s StoreWithContext) Get(key []byte) (Item, bool) { return s.store.Get(s.ctx, key) }
func (s StoreWithContext) Set(i Item) error            { return s.store.Set(s.ctx, i) }
func (s StoreWithContext) Delete(key []byte) error     { return s.store.Delete(s.ctx, key) }
func (s StoreWithContext) PrefixScan(prefix []byte) ([]Item, error) {
	return s.store.PrefixScan(s.ctx, prefix)
}
