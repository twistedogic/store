package bolt

import (
	"bytes"
	"context"
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/twistedogic/store"
)

type Store struct {
	bucketName string
	db         *bolt.DB
}

func New(file string) (store.Store, error) {
	db, err := bolt.Open(file, 0666, nil)
	if err != nil {
		return nil, err
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(file))
		return err
	}); err != nil {
		return nil, err
	}
	return Store{
		bucketName: file,
		db:         db,
	}, nil
}

func (s Store) Get(ctx context.Context, key []byte) (store.Item, bool) {
	item := store.Item{Key: key}
	err := s.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(s.bucketName)).Get(key)
		if len(v) == 0 {
			return fmt.Errorf("key %s not found", key)
		}
		item.Data = v
		return nil
	})
	return item, err == nil
}

func (s Store) Set(ctx context.Context, i store.Item) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(s.bucketName)).Put([]byte(i.Key), i.Data)
	})
}

func (s Store) Delete(ctx context.Context, key []byte) error {
	if _, exist := s.Get(ctx, key); !exist {
		return fmt.Errorf("key %s not found", key)
	}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(s.bucketName)).Delete(key)
	}); err != nil {
		return err
	}
	return s.db.Sync()
}

func (s Store) PrefixScan(ctx context.Context, prefix []byte) ([]store.Item, error) {
	out := make([]store.Item, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(s.bucketName)).ForEach(func(k, v []byte) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if bytes.HasPrefix(k, prefix) {
					out = append(out, store.Item{Key: k, Data: v})
				}
			}
			return nil
		})
	})
	return out, err
}
