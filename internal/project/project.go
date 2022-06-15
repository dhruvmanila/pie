package project

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type Project struct {
	Name    string
	Path    string
	VenvDir string
}

func New() (*Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	_, name := filepath.Split(cwd)

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
