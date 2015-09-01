package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestNewRepository(t *testing.T) {
	_, err := NewRepository("dummy")
	if err == nil {
		t.Error("NewRepository must returns error if a passed path is not a git directory")
	}

	path, err := createGitRepository()
	if err != nil {
		t.Fatal("failed to create a test git repository", err)
	}
	defer os.RemoveAll(path)

	if _, err := NewRepository(path); err != nil {
		t.Error("failed to instantiate repository")
	}
}

func TestRepository_Tags(t *testing.T) {
	path, err := createGitRepository()
	if err != nil {
		t.Fatal("failed to create a test git repository")
	}
	defer os.RemoveAll(path)

	repo, _ := NewRepository(path)
	tags, _ := repo.Tags()
	want := []string{"add_a"}
	for i, tag := range tags {
		if tag != want[i] {
			t.Errorf("repository.Tag returned a wrong tag. get %s, want %s", tag, want[i])
		}
	}
}

func TestRepository_Checkout(t *testing.T) {
	path, err := createGitRepository()
	if err != nil {
		t.Fatal("failed to create a test git repository")
	}
	defer os.RemoveAll(path)

	repo, _ := NewRepository(path)
	repo.Checkout("A")
	cmd := exec.Command("git", "symbolic-ref", "HEAD")
	cmd.Dir = path
	out, _ := cmd.CombinedOutput()
	ref := strings.TrimSpace(string(out))
	want := "refs/heads/A"
	if ref != want {
		t.Errorf("a wrong symbolic ref. got %s, want %s", ref, want)
	}
}

func TestRepository_Hash(t *testing.T) {
	path, err := createGitRepository()
	if err != nil {
		t.Fatal("failed to create a test git repository")
	}
	defer os.RemoveAll(path)

	repo, _ := NewRepository(path)
	got, _ := repo.Hash("master")

	cmd := exec.Command("git", "rev-parse", "master")
	cmd.Dir = path
	out, _ := cmd.CombinedOutput()
	want := strings.TrimSpace(string(out))

	if got != want {
		t.Errorf("repository.Hash() returned a wrong hash. got %s, wnat %s", got, want)
	}
}

func TestRepository_HeadHash(t *testing.T) {
	path, err := createGitRepository()
	if err != nil {
		t.Fatal("failed to create a test git repository")
	}
	defer os.RemoveAll(path)

	repo, _ := NewRepository(path)
	got, _ := repo.HeadHash()

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = path
	out, _ := cmd.CombinedOutput()
	want := strings.TrimSpace(string(out))

	if got != want {
		t.Errorf("repository.HeadHash() returned a wrong hash. got %s, wnat %s", got, want)
	}
}

func createGitRepository() (string, error) {
	dir, err := ioutil.TempDir("", "test-git")
	if err != nil {
		return "", err
	}

	cmds := []*exec.Cmd{
		exec.Command("git", "init"),
		exec.Command("git", "config", "user.name", "foo"),
		exec.Command("git", "config", "user.email", "foo@bar"),
		exec.Command("touch", "a"),
		exec.Command("git", "add", "a"),
		exec.Command("git", "commit", "-m", "add a"),
		exec.Command("git", "tag", "add_a"),
		exec.Command("git", "branch", "A"),
		exec.Command("touch", "b"),
		exec.Command("git", "add", "b"),
		exec.Command("git", "commit", "-m", "add b"),
	}

	for _, c := range cmds {
		c.Dir = dir
		if b, err := c.CombinedOutput(); err != nil {
			fmt.Println(string(b))
			return "", err
		}
	}

	return dir, nil
}
