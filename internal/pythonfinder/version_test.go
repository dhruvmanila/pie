package pythonfinder

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	pep440Version "github.com/aquasecurity/go-pep440-version"
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

func TestIsFinalRelease(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{
			version: "3.11.0",
			want:    true,
		},
		{
			version: "3.11.0a1",
			want:    false,
		},
		{
			version: "3.11.0.post1",
			want:    false,
		},
		{
			version: "3.11.0.dev1",
			want:    false,
		},
		{
			version: "3.11.0+local",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			v := pep440Version.MustParse(tt.version)
			got := isFinalRelease(&v)
			if got != tt.want {
				t.Errorf("isFinalRelease(%q) = %v, want %v", v, got, tt.want)
			}
		})
	}
}

func TestIsCompleteVersion(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{
			version: "3",
			want:    false,
		},
		{
			version: "3.11",
			want:    false,
		},
		{
			version: "3.11.0",
			want:    true,
		},
		{
			version: "3.11.0a1",
			want:    false,
		},
		{
			version: "3.11.0.post1",
			want:    false,
		},
		{
			version: "3.11.0.dev1",
			want:    false,
		},
		{
			version: "3.11.0+local",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			v := pep440Version.MustParse(tt.version)
			got := isCompleteVersion(&v)
			if got != tt.want {
				t.Errorf("isCompleteVersion(%q) = %v, want %v", v, got, tt.want)
			}
		})
	}
}

func TestGetGlobVersion(t *testing.T) {
	tests := []struct {
		version string
		want    string
	}{
		{
			version: "3",
			want:    "3.*",
		},
		{
			version: "3.11",
			want:    "3.11.*",
		},
		{
			version: "3.11.0",
			want:    "3.11.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			v := pep440Version.MustParse(tt.version)
			got := getGlobVersion(&v)
			if got != tt.want {
				t.Errorf("getGlobVersion(%q) = %v, want %v", v, got, tt.want)
			}
		})
	}
}
