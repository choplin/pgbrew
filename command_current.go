package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func DoCurrent(c *cli.Context) {
	args := c.Args()

	unset := c.Bool("unset")

	if unset {
		unsetCurrentVersion()
	} else if len(args) == 0 {
		showCurrentVersion()
	} else if len(args) == 1 {
		setCurrentVersion(args[0])
	} else {
		showHelpAndExit(c, "too many argument")
	}
}

func showCurrentVersion() {
	if !exists(currentLink) {
		fmt.Println("current version is not set")
		os.Exit(0)
	}

	if !isSymLink(currentLink) {
		log.WithField("path", currentLink).Fatalf("current path is not symbolic link")
	}

	path, err := os.Readlink(currentLink)
	if err != nil {
		log.WithField("err", err).Fatal("failed to read link")
	}
	fmt.Println(filepath.Base(path))
}

func setCurrentVersion(name string) {
	if exists(currentLink) {
		removeCurrentLink()
	}

	version, err := NewVersion(name)
	if err != nil {
		log.WithField("err", err).Fatal("failed to get version")
	}

	log.WithField("version", version.Name).Info("set a current version")
	log.WithFields(log.Fields{
		"link": currentLink,
		"to":   version.Path(),
	}).Debug("create a symbolic link")

	if err := os.Symlink(version.Path(), currentLink); err != nil {
		log.WithField("err", err).Fatal("failed to create a symbolic link")
	}
}

func unsetCurrentVersion() {
	if !exists(currentLink) {
		fmt.Println("current version is not set")
		os.Exit(0)
	}

	if !isSymLink(currentLink) {
		log.WithField("path", currentLink).Fatalf("current path is not symbolic link")
	}

	log.Info("unset a current version")
	removeCurrentLink()
}

func CurrentCompletion(c *cli.Context) {
	versions := AllVersions()
	for _, v := range versions {
		fmt.Println(v.Name)
	}
}

func isSymLink(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		log.WithField("err", err).Fatal("failed to get stat of file")
	}
	return fi.Mode()&os.ModeSymlink != 0
}

func removeCurrentLink() {
	log.WithField("path", currentLink).Debug("remove a current symbolic link")
	if err := os.Remove(currentLink); err != nil {
		log.WithField("err", err).Fatal("failed to remove an exisiting symbolic link")
	}
}
