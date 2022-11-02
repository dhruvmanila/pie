package pythonfinder

import (
	"os"
	"strings"
)

// pathProvider is a Provider that finds Python executables in the PATH
// environment variable.
type pathProvider struct {
	paths []string
}

// newPathProvider returns a new pathProvider.
func newPathProvider() *pathProvider {
	return &pathProvider{
		paths: strings.Split(os.Getenv("PATH"), string(os.PathListSeparator)),
	}
}

func (p *pathProvider) Executables() ([]string, error) {
	var executables []string
	for _, path := range p.paths {
		execs, err := execsInPath(path)
		if err != nil {
			return nil, err
		}
		executables = append(executables, execs...)
	}
	return executables, nil
}
