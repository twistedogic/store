package store

import (
	"context"
	"errors"
	"io"
)

var KeyNotFound = errors.New("key not found")

type Getter interface {
	// Get writes value of key into w if key exists
	// exist is true if key exists
	Get(ctx context.Context, key string, w io.Writer) (err error)
}

type Setter interface {
	Set(context.Context, string, io.Reader) error
}

type Deleter interface {
	Delete(context.Context, string) error
}

type PrefixScanner interface {
	PrefixScan(context.Context, string) ([]string, error)
}

type Store interface {
	Getter
	Setter
	Deleter
	PrefixScanner
}
