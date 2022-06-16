package project

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
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

// New creates a new project for the current directory.
func New() (*Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		return nil, err
	}
	_, name := filepath.Split(path)

	hash, err := hashPath(cwd)
	if err != nil {
		return nil, err
	}

	dataDir, err := xdg.DataFile("pyvenv/")
	if err != nil {
		return nil, err
	}

	venvName := fmt.Sprintf("%s-%s", name, hash[:8])
	return &Project{
		Name:    name,
		Path:    cwd,
		VenvDir: filepath.Join(dataDir, venvName),
	}, err
}

// WriteProjectFile associates the project directory with the virtual
// environment. This is done by writing the absolute path to the project
// directory in a ".project" file inside the virtual environment directory.
func (p *Project) WriteProjectFile() error {
	projectFilePath := filepath.Join(p.VenvDir, ".project")
	f, err := os.OpenFile(projectFilePath, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(p.Path)
	if err != nil {
		return err
	}

	return nil
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
