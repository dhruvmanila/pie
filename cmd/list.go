package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List out all the managed virtualenvs",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		dataDir, err := xdg.DataFile("pyvenv/")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Root directory: %s\n", green.Sprint(dataDir))

		venvs, err := getVenvs(dataDir)
		if err != nil {
			log.Fatal(err)
		}

		_, currentVenv := filepath.Split(os.Getenv("VIRTUAL_ENV"))
		for _, venv := range venvs {
			if currentVenv == venv.Name {
				bold.Print("* " + venv.Name)
			} else {
				bold.Print("  " + venv.Name)
			}
			yellowBold.Printf(" (%s)", venv.Version)
			faint.Printf(" (%s)\n", venv.Project)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// VirtualEnv contains information regarding a single virtual environment
// managed by the tool.
type VirtualEnv struct {
	// Name is the virtual environment name.
	Name string

	// Path is the absolute path to the virtual environment directory.
	Path string

	// Project is the absolute path to the project this virtual environment
	// belongs to.
	Project string

	// Version is the Python version this environment was created from.
	Version string
}

// getVenvs returns information regarding all the managed virtual environments.
func getVenvs(dataDir string) ([]*VirtualEnv, error) {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}

	var venvs []*VirtualEnv
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		venvName := entry.Name()
		venvDir := filepath.Join(dataDir, venvName)

		projectPath, err := readProjectFile(venvDir)
		if err != nil {
			return nil, err
		}

		pythonVersion, err := getPythonVersionFromConfig(venvDir)
		if err != nil {
			return nil, err
		}

		venvs = append(venvs, &VirtualEnv{
			Name:    venvName,
			Path:    venvDir,
			Project: projectPath,
			Version: pythonVersion,
		})
	}

	return venvs, nil
}

// getPythonVersionFromConfig returns the Python version used to create
// the virtual environment.
//
// The version string is read from the environment config file present
// in the given environment directory.
func getPythonVersionFromConfig(venvDir string) (string, error) {
	pyvenvPath := filepath.Join(venvDir, "pyvenv.cfg")

	file, err := os.Open(pyvenvPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pairs := strings.Split(scanner.Text(), "=")
		if strings.TrimSpace(pairs[0]) == "version" {
			return strings.TrimSpace(pairs[1]), nil
		}
	}

	return "", fmt.Errorf("%q: venv config file does not contain 'version' key", pyvenvPath)
}
