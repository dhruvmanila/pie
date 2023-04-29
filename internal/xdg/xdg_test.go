package xdg_test

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/dhruvmanila/pie/internal/xdg"
)

func TestDataDir(t *testing.T) {
	if xdg.DataDir == "" {
		t.Fatal("DataDir is empty")
	}

	if fi, err := os.Stat(xdg.DataDir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			t.Errorf("%q: directory does not exist", xdg.DataDir)
		}
		t.Errorf("%q: stat error = %v", xdg.DataDir, err)
	} else if !fi.IsDir() {
		t.Errorf("%q: not a directory", xdg.DataDir)
	}
}
