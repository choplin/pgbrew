package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/choplin/pgenv/git"
)

const (
	gitUrl = "git://git.postgresql.org/git/postgresql"
)

func DoClone(c *cli.Context) {
	log.Info("clone postgresql git repository")
	if err := git.Clone(localRepository, gitUrl, c.Args()); err != nil {
		log.WithField("err", err).Fatal("failed to clone git reporitory")
	}
}
