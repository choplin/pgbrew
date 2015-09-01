package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetCurrentVersion(t *testing.T) {
	app := makeTestEnv()
	target := "9.4.4"
	app.Run([]string{"pgbrew", "current", target})

	fi, err := os.Lstat(currentLink)
	if err != nil {
		t.Fatalf("failed to retrieve a FileInfo of %s", currentLink)
	}

	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("%s is not a symbolic link", currentLink)
	}

	path, err := os.Readlink(currentLink)
	if err != nil {
		t.Fatal("failed to resolve a symbolic link")
	}

	expected := filepath.Join(installBase, target)
	if path != expected {
		t.Errorf("a wrong link of current version. got: %s, want: %s", expected, path)
	}
}

func ExampleCurrentCommand_showCurrentLink() {
	app := makeTestEnv()

	target := "9.4.4"
	app.Run([]string{"pgbrew", "current", target})

	app.Run([]string{"pgbrew", "current"})
	// Output: 9.4.4
}

func TestUnsetCurrentVersion(t *testing.T) {
	app := makeTestEnv()

	target := "9.4.4"
	app.Run([]string{"pgbrew", "current", target})
	if !exists(currentLink) {
		t.Fatal("failed to set a current version")
	}

	app.Run([]string{"pgbrew", "current", "-u"})
	if exists(currentLink) {
		t.Error("failed to unset a current version")
	}
}
