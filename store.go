package store

import (
	"context"
)

type Item struct {
	Key  string
	Data []byte
}

type Getter interface {
	Exists(context.Context, string) (bool, error)
	Get(context.Context, string) (Item, error)
}

type Setter interface {
	Set(context.Context, Item) error
}

type Deleter interface {
	Delete(context.Context, string) error
}

type PrefixScanner interface {
	PrefixScan(context.Context, string) ([]Item, error)
}

type Store interface {
	Getter
	Setter
	Deleter
	PrefixScanner
}
