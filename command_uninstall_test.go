package main

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestUninstallCommand(t *testing.T) {
	target := "test"

	src := filepath.Join(baseDir.installDir(), "9.4.4")
	dest := filepath.Join(baseDir.installDir(), target)
	cmd := exec.Command("cp", "-prf", src, dest)
	if err := cmd.Run(); err != nil {
		t.Fatal("failed to run cp command")
	}

	app := makeTestEnv()
	app.Run([]string{"pgenv", "uninstall", target})

	if exists(dest) {
		t.Errorf("found uninstalled version %s", dest)
	}
}
