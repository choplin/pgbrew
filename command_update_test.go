package main

import (
	"testing"

	"github.com/choplin/pgbrew/git"
)

func TestUpdateCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skip a test for update command")
	}

	repo, err := git.NewRepository(localRepository)
	if err != nil {
		t.Fatal("failed to initialize a reporitory")
	}

	hash, err := repo.Hash("origin/HEAD")
	if err != nil {
		t.Fatal("failed to get a hash of origin/HEAD")
	}

	app := makeTestEnv()
	app.Run([]string{"pgbrew", "update"})

	updated, err := repo.Hash("origin/HEAD")
	if err != nil {
		t.Fatal("failed to get a hash of origin/HEAD")
	}

	if hash == updated {
		t.Error("local repository is not updated")
	}
}
