package xdg

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const appName = "pyvenv"

// DataDir defines the directory where pyvenv stores all the virtual
// environments.
var DataDir string

func init() {
	DataDir = filepath.Join(xdg.DataHome, appName)
	if _, err := os.Stat(DataDir); errors.Is(err, fs.ErrNotExist) {
		if err := os.Mkdir(DataDir, 0o755); err != nil {
			log.Fatal(err)
		}
	}
}
