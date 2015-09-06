package main

import (
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func DoInitdb(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
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

	options := c.String("options")

	initdbArgs := []string{"-D", path}
	if options != "" {
		for _, o := range strings.Split(options, " ") {
			initdbArgs = append(initdbArgs, o)
		}
	}

	log.WithField("options", initdbArgs).Info("initdb")
	if err := pg.Initdb(initdbArgs); err != nil {
		log.WithField("err", err).Fatal("failed to execute initdb")
	}

	writeClusterVersionFile(name, pg)
}

func InitdbCompletion(c *cli.Context) {
	args := c.Args()

	if len(args) == 0 {
		versions := AllVersions()
		for _, v := range versions {
			fmt.Println(v.Name)
		}
	}
}

func writeClusterVersionFile(name string, pg *Postgres) {
	cluster := &Cluster{
		pg:   pg,
		name: name,
	}

	log.WithFields(log.Fields{
		"version name": pg.Version().Name,
	}).Info("write a cluster version file")

	if err := cluster.WriteVersionFile(); err != nil {
		log.WithField("err", err).Fatal("failed to write a cluster version file")
	}
}
