package main

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestUninstallCommand(t *testing.T) {
	target := "test"

	src := filepath.Join(installBase, "9.4.4")
	dest := filepath.Join(installBase, target)
	cmd := exec.Command("cp", "-prf", src, dest)
	if err := cmd.Run(); err != nil {
		t.Fatal("failed to run cp command")
	}

	app := makeTestEnv()
	app.Run([]string{"pgbrew", "uninstall", target})

	if exists(dest) {
		t.Errorf("found uninstalled version %s", dest)
	}
}
