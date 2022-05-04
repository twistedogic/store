package mem

import (
	"testing"

	"github.com/twistedogic/store"
	"github.com/twistedogic/store/testutil"
)

func Test_Store(t *testing.T) {
	factory := func(t *testing.T) store.Store {
		s, err := New()
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	testutil.RunStoreTest(t, factory)
}
