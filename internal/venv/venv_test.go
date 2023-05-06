package venv

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dhruvmanila/pie/internal/xdg"
)

var testdataDir string

func init() {
	dir, err := filepath.Abs("testdata")
	if err != nil {
		panic(err)
	}
	testdataDir = dir
}

func TestNames(t *testing.T) {
	originalDataDir := xdg.DataDir
	xdg.DataDir = testdataDir
	t.Cleanup(func() {
		xdg.DataDir = originalDataDir
	})

	venvNames, err := Names()
	if err != nil {
		t.Fatalf("Names() error = %v, want nil", err)
	}
	if len(venvNames) != 2 {
		t.Errorf("Names() len = %v, want 2", len(venvNames))
	}

	want := []string{"venv1", "venv2"}
	if !reflect.DeepEqual(venvNames, want) {
		t.Errorf("Names() = %v, want %s", venvNames, want)
	}
}

func TestProjectPath(t *testing.T) {
	originalDataDir := xdg.DataDir
	xdg.DataDir = testdataDir
	t.Cleanup(func() {
		xdg.DataDir = originalDataDir
	})

	got, err := ProjectPath("venv1")
	if err != nil {
		t.Fatalf("ProjectPath() error = %v, want nil", err)
	}

	want := "/home/user/project1"
	if got != want {
		t.Errorf("ProjectPath() = %v, want %s", got, want)
	}
}

func TestPythonVersion(t *testing.T) {
	originalDataDir := xdg.DataDir
	xdg.DataDir = testdataDir
	t.Cleanup(func() {
		xdg.DataDir = originalDataDir
	})

	got, err := PythonVersion("venv1")
	if err != nil {
		t.Fatalf("PythonVersion() error = %v, want nil", err)
	}

	want := "3.11.0"
	if got != want {
		t.Errorf("PythonVersion() = %v, want %s", got, want)
	}
}
