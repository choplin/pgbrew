package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isEmpty(dirname string) bool {
	fs, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(fmt.Sprintf("failed to read directory %s", err.Error()))
	}

	return len(fs) == 0
}

func IsGitRepository(dirname string) bool {
	return exists(filepath.Join(dirname, ".git"))
}
