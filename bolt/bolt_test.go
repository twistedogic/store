package bolt

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/twistedogic/store"
	"github.com/twistedogic/store/testutil"
)

func Test_Store(t *testing.T) {
	dir, err := ioutil.TempDir("", "bolt_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	factory := func(t *testing.T) store.Store {
		s, err := New(filepath.Join(dir, "bolt_db"))
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	testutil.RunStoreTest(t, factory)
}
