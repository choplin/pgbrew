package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/choplin/pgenv/git"
)

func DoAvailable(c *cli.Context) {
	repo, err := git.NewRepository(localRepository)
	if err != nil {
		log.WithField("err", err).Fatal("failed to initialize local reporitory")
	}

	tags, err := repo.Tags()
	if err != nil {
		log.WithField("err", err).Fatal("failed to get tags")
	}

	fmt.Println("Available versions:")
	for _, tag := range tags {
		fmt.Printf("\t%s\n", tag)
	}
}
