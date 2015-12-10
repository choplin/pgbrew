package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/choplin/pgenv/git"
	"github.com/codegangsta/cli"
)

const gitURL = "git://git.postgresql.org/git/postgresql"

func DoClone(c *cli.Context) {
	localRepository := config.RepositoryPath
	if localRepository != "" {
		if git.IsGitRepository(config.RepositoryPath) {
			log.Fatalf("local repository %s already exists", localRepository)
		} else {
			log.Fatalf("local repository %s already exists, but seems not be a git repository", localRepository)
		}
	}

	config.RepositoryPath = baseDir.defaultLocalRepository()
	log.WithFields(log.Fields{
		"config": config,
		"path":   configFilePath,
	}).Debug("write config file")
	if err := config.Write(configFilePath); err != nil {
		log.WithField("err", err).Fatal("failed to write a config file")
	}

	log.Info("clone postgresql git repository")
	if err := git.Clone(localRepository, gitURL, c.Args()); err != nil {
		log.WithField("err", err).Fatal("failed to clone git reporitory")
	}
}
