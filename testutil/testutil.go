package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/twistedogic/store"
)

type action string

const (
	GET    action = "GET"
	SET    action = "SET"
	DELETE action = "DELETE"
	SCAN   action = "SCAN"
	EXIST  action = "EXIST"
)

type factory func(t *testing.T) store.Store

type testOps struct {
	item            store.Item
	hasError, exist bool
	want            []store.Item
	action          action
}

func (o testOps) String() string {
	return fmt.Sprintf("%s %v", o.action, o.item)
}

func (o testOps) check(t *testing.T, err error) {
	if (err != nil) != o.hasError {
		t.Fatalf("%s: expect err %v got %v", o, o.hasError, err)
	}
}

func (o testOps) Run(t *testing.T, s store.Store) {
	ctx := context.TODO()
	key := o.item.Key
	var err error
	got, exist := make([]store.Item, 0), false
	switch o.action {
	case GET:
		g, ok := s.Get(ctx, key)
		if !ok {
			err = fmt.Errorf("key %q not found", key)
		}
		got = append(got, g)
	case SET:
		err = s.Set(ctx, o.item)
	case DELETE:
		err = s.Delete(ctx, key)
	case SCAN:
		got, err = s.PrefixScan(ctx, key)
	case EXIST:
		if _, exist = s.Get(ctx, key); exist != o.exist {
			t.Fatalf("%s: exist want: %v, got: %v", o, o.exist, exist)
		}
	}
	o.check(t, err)
	if err != nil {
		return
	}
	opt := cmpopts.SortSlices(func(a, b store.Item) bool {
		return string(a.Key) > string(b.Key)
	})
	if diff := cmp.Diff(o.want, got, opt); diff != "" {
		t.Fatalf("%s:\n%s", o, diff)
	}
}

func RunStoreTest(t *testing.T, f factory) {
	s := f(t)
	steps := []testOps{
		{
			item:   store.Item{Key: []byte("a"), Data: []byte("a")},
			want:   []store.Item{},
			action: SET,
		},
		{
			item:   store.Item{Key: []byte("b"), Data: []byte("b")},
			want:   []store.Item{},
			action: SET,
		},
		{
			item: store.Item{Key: []byte("a")},
			want: []store.Item{
				{Key: []byte("a"), Data: []byte("a")},
			},
			action: GET,
		},
		{
			item: store.Item{Key: []byte("")},
			want: []store.Item{
				{Key: []byte("a"), Data: []byte("a")},
				{Key: []byte("b"), Data: []byte("b")},
			},
			action: SCAN,
		},
		{
			item:   store.Item{Key: []byte("c")},
			want:   []store.Item{},
			exist:  false,
			action: EXIST,
		},
		{
			item:     store.Item{Key: []byte("c")},
			want:     []store.Item{},
			hasError: true,
			action:   GET,
		},
		{
			item:     store.Item{Key: []byte("c")},
			hasError: true,
			action:   DELETE,
		},
		{
			item:   store.Item{Key: []byte("a")},
			want:   []store.Item{},
			action: DELETE,
		},
		{
			item:     store.Item{Key: []byte("a")},
			want:     []store.Item{},
			hasError: true,
			action:   DELETE,
		},
		{
			item:   store.Item{Key: []byte("a")},
			want:   []store.Item{},
			action: SCAN,
		},
		{
			item: store.Item{Key: []byte("")},
			want: []store.Item{
				{Key: []byte("b"), Data: []byte("b")},
			},
			action: SCAN,
		},
		{
			item: store.Item{},
			want: []store.Item{
				{Key: []byte("b"), Data: []byte("b")},
			},
			action: SCAN,
		},
	}
	for _, step := range steps {
		step.Run(t, s)
	}
}
