package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/choplin/pgenv/git"
	"github.com/codegangsta/cli"
)

func DoInit(c *cli.Context) {
	if exists(configFilePath) {
		fmt.Println("already initilized")
		os.Exit(0)
	}

	log.Info("initialize pgenv")

	basePath := c.String("base-path")
	if basePath == "" {
		basePath = filepath.Join(homeDirectory, ".pgenv")
	}
	if exists(basePath) {
		log.Fatalf("base-path %s already exists", basePath)
	}

	doClone := true
	repositoryPath := c.String("repository-path")
	if repositoryPath != "" {
		doClone = false
	}

	dirs := []string{
		filepath.Dir(configFilePath),
		basePath,
	}
	for _, d := range dirs {
		if !exists(d) {
			log.WithField("path", d).Debug("create a directory")
			if err := os.Mkdir(d, 0755); err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"path": d,
				}).Fatal("failed to make a directory")
			}
		}
	}

	config := &Config{
		BasePath: basePath,
	}
	if repositoryPath != "" {
		if !git.IsGitRepository(repositoryPath) {
			log.Fatalf("repository-path %s is not a git repository", repositoryPath)
		}
		config.RepositoryPath = repositoryPath
	}

	log.WithFields(log.Fields{
		"config": config,
		"path":   configFilePath,
	}).Debug("write config file")
	if err := config.Write(configFilePath); err != nil {
		log.WithField("err", err).Fatal("failed to write a config file")
	}

	if doClone {
		println("clone")
	}
}
