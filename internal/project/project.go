package project

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhruvmanila/pie/internal/pathutil"
	"github.com/dhruvmanila/pie/internal/xdg"
)

// Project contains information regarding a specific project for which the
// virtual environment is being managed.
type Project struct {
	// Name is the project name.
	Name string

	// Path is the absolute path to the project directory.
	Path string

	// VenvDir is the absolute path to the virtual environment directory
	// for this project.
	VenvDir string
}

// New creates a new project for the given path after evaluating all the
// symlinks in the path. The given path must be an absolute path.
func New(path string) (*Project, error) {
	var err error
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return nil, err
	}

	_, name := filepath.Split(path)

	hash, err := hashPath(path)
	if err != nil {
		return nil, err
	}

	venvName := fmt.Sprintf("%s-%s", name, hash[:8])
	return &Project{
		Name:    name,
		Path:    path,
		VenvDir: filepath.Join(xdg.DataDir, venvName),
	}, err
}

// NewFromWd creates a new project from the current working directory.
func NewFromWd() (*Project, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return New(dir)
}

// Current returns the project associated with the current working directory,
// nil if there is none.
//
// The search starts from the current working directory and goes up the
// directory tree until a project is found or the root directory is reached.
func Current() (*Project, error) {
	p, err := NewFromWd()
	if err != nil {
		return nil, err
	}

	// root is the system root directory. For windows, it will be "C:\" while
	// for other systems it will be "/".
	root := filepath.VolumeName(p.Path) + string(os.PathSeparator)

	for p.Path != root {
		if pathutil.IsDir(p.VenvDir) {
			return p, nil
		}

		p, err = New(filepath.Dir(p.Path))
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// WriteProjectFile associates the project directory with the virtual
// environment. This is done by writing the absolute path to the project
// directory in a ".project" file inside the virtual environment directory.
func (p *Project) WriteProjectFile() error {
	return os.WriteFile(filepath.Join(p.VenvDir, ".project"), []byte(p.Path), 0o644)
}

// hashPath returns the hash value of the given path string. It uses the SHA 256
// algorithm to create the hash value.
func hashPath(path string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(path))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
