package pythonfinder

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

// Fake exec command helper
var testCaseName string

func TestPythonExecutable(t *testing.T) {
	// Setup fake exec command
	testCaseName = "TestGetVersionInfo"
	execCommand = fakeExecCommand
	defer func() {
		execCommand = exec.Command
	}()

	executable := "/bin/python"
	got, err := newPythonExecutable(executable)
	if err != nil {
		t.Fatalf("newPythonVersion(%q) unexpected error = %q", executable, err)
	}

	want := "3.11.0"
	if !reflect.DeepEqual(got.Version.String(), want) {
		t.Errorf("newPythonExecutable(%q).Version = %q, want %q", executable, got.Version, want)
	}
	if !reflect.DeepEqual(got.Path, executable) {
		t.Errorf("newPythonExecutable(%q).Path = %q, want %q", executable, got.Path, executable)
	}
}

func fakeExecCommand(name string, arg ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", name}
	cs = append(cs, arg...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", "GO_TEST_CASE_NAME=" + testCaseName}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	cmd, args := os.Args[3], os.Args[4:]

	if !strings.HasSuffix(cmd, "python") {
		fmt.Fprintf(os.Stderr, "command not found: %q", cmd)
		os.Exit(1)
	}

	if len(args) != 1 || args[0] != "--version" {
		fmt.Fprintf(os.Stderr, "invalid arguments: %q", args)
		os.Exit(1)
	}

	switch os.Getenv("GO_TEST_CASE_NAME") {
	case "TestGetVersionInfo":
		fmt.Fprintln(os.Stdout, "Python 3.11.0")
	}

	os.Exit(0)
}
