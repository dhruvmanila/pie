package pythonfinder

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestPathProvider(t *testing.T) {
	want := []string{
		"/bin",
		"/usr/bin",
		"/usr/local/bin",
	}
	t.Setenv("PATH", strings.Join(want, string(os.PathListSeparator)))

	p := newPathProvider()
	if !reflect.DeepEqual(p.paths, want) {
		t.Errorf("got %q, want %q", p.paths, want)
	}
}
