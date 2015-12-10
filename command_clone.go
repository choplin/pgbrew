package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/choplin/pgenv/git"
	"github.com/codegangsta/cli"
)

const gitURL = "git://git.postgresql.org/git/postgresql"

func DoClone(c *cli.Context) {
	if localRepository != "" {
		if git.IsGitRepository(localRepository) {
			log.Fatalf("local repository %s already exists", localRepository)
		}
	}

	config.RepositoryPath = localRepository
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
