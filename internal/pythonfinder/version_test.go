package pythonfinder

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

// Fake exec command helper
var testCaseName string

func TestVersionInfo(t *testing.T) {
	v1 := &VersionInfo{Major: 3, Minor: 11, Patch: 0}
	v2 := &VersionInfo{Major: 3, Minor: 11, Patch: 0}

	if !v1.Matches(v2) {
		t.Errorf("(%#v).Matches(%#v) = false, want true", v1, v2)
	}

	got := v1.String()
	want := "3.11.0"
	if got != want {
		t.Errorf("(%#v).String() = %q, want %q", v1, got, want)
	}
}

func TestGetVersionInfo_Valid(t *testing.T) {
	// Setup fake exec command
	testCaseName = "TestGetVersionInfoValid"
	execCommand = fakeExecCommand
	defer func() {
		execCommand = exec.Command
	}()

	executable := "/bin/python"
	got, err := getVersionInfo(executable)
	if err != nil {
		t.Fatalf("getVersionInfo(%q) unexpected error = %v", executable, err)
	}

	want := &VersionInfo{Major: 3, Minor: 11, Patch: 0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("getVersionInfo(%q) = %#v, want %#v", executable, got, want)
	}
}

func TestGetVersionInfo_Invalid(t *testing.T) {
	// Setup fake exec command
	testCaseName = "TestGetVersionInfoInvalid"
	execCommand = fakeExecCommand
	defer func() {
		execCommand = exec.Command
	}()

	executable := "/bin/python"
	_, err := getVersionInfo(executable)
	if !errors.Is(err, ErrInvalidVersion) {
		t.Fatalf("getVersionInfo(%q) error = %q, want = %q", executable, err, ErrInvalidVersion)
	}
}

func TestPythonVersion(t *testing.T) {
	// Setup fake exec command
	testCaseName = "TestGetVersionInfoValid"
	execCommand = fakeExecCommand
	defer func() {
		execCommand = exec.Command
	}()

	executable := "/bin/python"
	got, err := newPythonVersion(executable)
	if err != nil {
		t.Fatalf("newPythonVersion(%q) unexpected error = %q", executable, err)
	}

	want := &PythonVersion{
		Executable:  executable,
		VersionInfo: &VersionInfo{Major: 3, Minor: 11, Patch: 0},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("newPythonVersion(%q) = %#v, want %#v", executable, got, want)
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

	args := os.Args
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
	case "TestGetVersionInfoValid":
		fmt.Fprintln(os.Stdout, "Python 3.11.0")
	case "TestGetVersionInfoInvalid":
		fmt.Fprintln(os.Stdout, "Python 3.12.0a1")
	}

	os.Exit(0)
}
