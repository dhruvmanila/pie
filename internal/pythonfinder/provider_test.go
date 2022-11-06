package pythonfinder

import (
	"io/fs"
	"runtime"
	"testing"
	"time"
)

func TestLooksLikePython(t *testing.T) {
	var tests map[string]bool

	switch runtime.GOOS {
	case "windows":
		tests = map[string]bool{
			"python.exe":         true,
			"python3.exe":        true,
			"python39.exe":       true,
			"python310.exe":      true,
			"python-build":       false,
			"python-python3.exe": false,
		}
	default:
		tests = map[string]bool{
			"python":         true,
			"python3":        true,
			"python3.9":      true,
			"python3.10":     true,
			"python-build":   false,
			"python-python3": false,
		}
	}

	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			got := looksLikePython(name)
			if got != want {
				t.Errorf("looksLikePython(%q) = %v, want %v", name, got, want)
			}
		})
	}
}

type fakeFileInfo struct {
	dir      bool
	basename string
	modtime  time.Time
	mode     fs.FileMode
}

func (f *fakeFileInfo) Name() string       { return f.basename }
func (f *fakeFileInfo) Size() int64        { return 0 }
func (f *fakeFileInfo) Mode() fs.FileMode  { return f.mode }
func (f *fakeFileInfo) ModTime() time.Time { return f.modtime }
func (f *fakeFileInfo) IsDir() bool        { return f.dir }
func (f *fakeFileInfo) Sys() interface{}   { return nil }

func TestIsExecutable(t *testing.T) {
	var tests map[*fakeFileInfo]bool

	switch runtime.GOOS {
	case "windows":
		tests = map[*fakeFileInfo]bool{
			{basename: "python.exe"}: true,
			{basename: "python"}:     false,
		}
	default:
		tests = map[*fakeFileInfo]bool{
			{mode: fs.ModeDir}: false,
			{mode: 0o100}:      true,
			{mode: 0o110}:      true,
			{mode: 0o111}:      true,
		}
	}

	for info, want := range tests {
		t.Run(info.Mode().String(), func(t *testing.T) {
			got := isExecutable(info)
			if got != want {
				t.Errorf("isExecutable(%v) = %v, want %v", info, got, want)
			}
		})
	}
}
