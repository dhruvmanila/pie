package project

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dhruvmanila/pyvenv/internal/xdg"
)

var testdataDir string

func init() {
	dir, err := filepath.Abs("testdata")
	if err != nil {
		panic(err)
	}
	testdataDir = dir
}

// chdir changes the current working directory to the named directory,
// and then restore the original working directory at the end of the test.
func chdir(t *testing.T, dir string) {
	olddir, err := os.Getwd()
	if err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(olddir); err != nil {
			t.Fatalf("chdir to original working directory %s: %v", olddir, err)
		}
	})
}

// setupVenvDir creates a virtual environment directory for the given project
// and returns the path to the virtual environment directory. The venv directory
// will be removed at the end of the test.
func setupVenvDir(t *testing.T, dir string) string {
	hash, err := hashPath(dir)
	if err != nil {
		t.Fatalf("hashPath(%q) error = %v", dir, err)
	}

	_, name := filepath.Split(dir)
	venvDir := filepath.Join(xdg.DataDir, fmt.Sprintf("%s-%s", name, hash[:8]))
	if err = os.MkdirAll(venvDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", venvDir, err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(venvDir); err != nil {
			t.Fatalf("RemoveAll(%q) error = %v", venvDir, err)
		}
	})
	return venvDir
}

func verifyProject(t *testing.T, p *Project, dir string, newFunc string) {
	_, projectName := filepath.Split(dir)
	if !strings.HasPrefix(p.Name, projectName) {
		t.Errorf("%s(%q).Name = %q, want %s*", newFunc, dir, p.Name, projectName)
	}

	if p.Path != dir {
		t.Errorf("%s(%q).Path = %q, want %q", newFunc, dir, p.Path, dir)
	}

	want := filepath.Join(xdg.DataDir, projectName)
	if !strings.HasPrefix(p.VenvDir, want) {
		t.Errorf("%s(%q).VenvDir = %q, want %q*", newFunc, dir, p.VenvDir, want)
	}
}

func TestNewProject(t *testing.T) {
	alpha := filepath.Join(testdataDir, "alpha")
	p, err := New(alpha)
	if err != nil {
		t.Fatalf("New(%q) error = %v, want nil", alpha, err)
	}

	verifyProject(t, p, alpha, "New")

	beta := filepath.Join(testdataDir, "beta")
	p, err = New(beta)
	if err != nil {
		t.Fatalf("New(%q) error = %v, want nil", beta, err)
	}

	// beta is a symlink to alpha.
	verifyProject(t, p, alpha, "New")

	chdir(t, alpha)
	p, err = NewFromWd()
	if err != nil {
		t.Fatalf("NewFromWd(%q) error = %v, want nil", alpha, err)
	}

	verifyProject(t, p, alpha, "NewFromWd")
}

func TestCurrentProjectNoVenv(t *testing.T) {
	alpha := filepath.Join(testdataDir, "alpha")

	// Alpha does not have a virtual environment associated with it.
	chdir(t, alpha)

	p, err := Current()
	if err != nil {
		t.Errorf("Current() error = %v, want nil", err)
	}
	if p != nil {
		t.Errorf("Current() = %v, want nil", p)
	}
}

func TestCurrentProjectWithVenv(t *testing.T) {
	parent := filepath.Join(testdataDir, "parent")
	setupVenvDir(t, parent)
	chdir(t, parent)

	p, err := Current()
	if err != nil {
		t.Errorf("Current() error = %v, want nil", err)
	}
	if p == nil {
		t.Fatal("Current() = nil, want non-nil")
	}

	verifyProject(t, p, parent, "Current")

	// Child is a subdirectory of parent.
	child := filepath.Join(parent, "child")
	chdir(t, child)

	p, err = Current()
	if err != nil {
		t.Errorf("Current() error = %v, want nil", err)
	}
	if p == nil {
		t.Fatal("Current() = nil, want non-nil")
	}

	// The virtual environment is for the parent project, so verify that
	// it traverses up the directory tree to find the parent project.
	verifyProject(t, p, parent, "Current")
}

func TestWriteProjectFile(t *testing.T) {
	alpha := filepath.Join(testdataDir, "alpha")
	venvDir := setupVenvDir(t, alpha)

	p, err := New(alpha)
	if err != nil {
		t.Fatalf("New(%q) error = %v, want nil", alpha, err)
	}

	// Write the project file.
	if err = p.WriteProjectFile(); err != nil {
		t.Fatalf("WriteProjectFile() error = %v, want nil", err)
	}

	// Verify that the project file was written.
	projectFile := filepath.Join(venvDir, ".project")
	if _, err = os.Stat(projectFile); err != nil {
		t.Fatalf("Stat(%q) error = %v, want nil", projectFile, err)
	}

	// Verify that the project file contains the project path.
	b, err := os.ReadFile(projectFile)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v, want nil", projectFile, err)
	}
	if !bytes.Equal(b, []byte(alpha)) {
		t.Errorf("ReadFile(%q) = %q, want %q", projectFile, b, p.Path)
	}
}
