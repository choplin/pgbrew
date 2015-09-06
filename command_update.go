package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/choplin/pgbrew/git"
)

func DoUpdate(c *cli.Context) {
	log.Info("update a local git repository")
	repo, err := git.NewRepository(localRepository)
	if err != nil {
		log.WithField("err", err).Fatal("failed to initialize a reporitory")
	}

	if out, err := repo.Fetch(); err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
			"out": out,
		}).Fatal("failed to fetch")
	}
}
