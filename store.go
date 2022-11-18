package store

import (
	"context"
	"io"
)

type Getter interface {
	Get(context.Context, string, io.Writer) (bool, error)
}

type Setter interface {
	Set(context.Context, io.Reader, ...Metadata) error
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
