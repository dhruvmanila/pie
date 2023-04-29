package pathutil_test

import (
	"os"
	"testing"

	"github.com/dhruvmanila/pie/internal/pathutil"
)

func TestIsDir(t *testing.T) {
	tempdir := t.TempDir()
	if !pathutil.IsDir(tempdir) {
		t.Errorf("IsDir(%q) = false, want true", tempdir)
	}

	tempfile, err := os.CreateTemp(tempdir, "file")
	if err != nil {
		t.Fatal(err)
	}
	defer tempfile.Close()

	if pathutil.IsDir(tempfile.Name()) {
		t.Errorf("IsDir(%q) = true, want false", tempfile.Name())
	}
}
