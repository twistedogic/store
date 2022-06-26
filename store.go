package store

import (
	"context"
)

type Getter interface {
	Get(context.Context, []byte) (Item, bool)
}

type Setter interface {
	Set(context.Context, Item) error
}

type Deleter interface {
	Delete(context.Context, []byte) error
}

type PrefixScanner interface {
	PrefixScan(context.Context, []byte) ([]Item, error)
}

type Store interface {
	Getter
	Setter
	Deleter
	PrefixScanner
}
