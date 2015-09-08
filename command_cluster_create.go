package main

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func DoClusterCreate(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		showHelpAndExit(c, "<version> must be specified")
	}

	pg, err := NewPostgres(args[0])
	if err != nil {
		log.WithField("err", err).Fatal("a specified version is not installed")
	}

	name := c.String("name")
	if name == "" {
		name = pg.Version().Name
	}

	path := filepath.Join(clusterBase, name)
	if exists(path) {
		log.WithField("name", name).Fatal("cluster already exists")
	}

	initdbArgs := []string{"-D", path}
	if len(args) > 1 {
		for _, a := range args[1:] {
			initdbArgs = append(initdbArgs, a)
		}
	}

	log.WithField("initdb options", initdbArgs).Debug("initdb")
	if err := pg.Initdb(initdbArgs); err != nil {
		log.WithField("err", err).Fatal("failed to execute initdb")
	}

	writeClusterExtraInfoFile(name, pg)
}

func ClusterCreateCompletion(c *cli.Context) {
	args := c.Args()

	if len(args) == 0 {
		versions := AllVersions()
		for _, v := range versions {
			fmt.Println(v.Name)
		}
	}
}

func writeClusterExtraInfoFile(name string, pg *Postgres) {
	cluster := &Cluster{
		Pg:   pg,
		Name: name,
	}

	log.WithFields(log.Fields{
		"version name": pg.Version().Name,
	}).Debug("write a cluster extra info file")

	if err := cluster.WriteExtraInfoFile(); err != nil {
		log.WithField("err", err).Fatal("failed to write a cluster extra info file")
	}
}
