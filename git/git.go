package git

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/choplin/pgenv/util"
)

func Clone(path string, url string, options []string) error {
	if exists(path) {
		return errors.New("local reporitory already exists")
	}

	args := []string{"clone"}
	for _, o := range options {
		args = append(args, o)
	}
	args = append(args, url)
	args = append(args, path)

	cmd := exec.Command("git", args...)
	return util.RunCommandWithDebugLog(cmd)
}

type repository struct {
	path string
}

func NewRepository(path string) (*repository, error) {
	if !IsGitRepository(path) {
		return nil, errors.New("local reporitory is not initialized")
	}

	return &repository{
		path: path,
	}, nil
}

func (r *repository) Tags() ([]string, error) {
	cmd := r.gitCommand("show-ref", "--tags")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	tags := make([]string, 0)

	for _, l := range lines {
		if !strings.Contains(l, "refs/tags/") {
			continue
		}
		ref := strings.Split(l, " ")[1]
		tags = append(tags, strings.Split(ref, "/")[2])
	}

	return tags, nil
}

func (r *repository) Checkout(gitRef string) (string, error) {
	cmd := r.gitCommand("checkout", "-q", "-f", gitRef)
	if out, err := cmd.CombinedOutput(); err != nil {
		return string(out), err
	}
	return "", nil
}

func (r *repository) CheckoutWithWorkTree(gitRef string, worktree string) (string, error) {
	cmd := r.gitCommand("--work-tree", worktree, "checkout", "-q", "-f", gitRef)
	if out, err := cmd.CombinedOutput(); err != nil {
		return string(out), err
	}
	return "", nil
}

func (r *repository) Hash(ref string) (string, error) {
	cmd := r.gitCommand("rev-parse", "--verify", ref)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func (r *repository) HeadHash() (string, error) {
	return r.Hash("HEAD")
}

func (r *repository) gitCommand(commands ...string) *exec.Cmd {
	args := []string{"--git-dir", filepath.Join(r.path, ".git")}
	for _, c := range commands {
		args = append(args, c)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = r.path
	return cmd
}

func (r *repository) Fetch() (string, error) {
	cmd := r.gitCommand("fetch", "-q")
	if out, err := cmd.CombinedOutput(); err != nil {
		return string(out), err
	}
	return "", nil
}
