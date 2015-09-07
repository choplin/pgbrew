package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/choplin/pgbrew/util"
	"github.com/codegangsta/cli"
)

func DoClusterStart(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		showHelpAndExit(c, "<cluster> must be specified")
	}

	cluster, err := NewCluster(args[0])
	if err != nil {
		log.WithField("err", err).Fatal("failed to get a cluster")
	}

	if cluster.IsRunning() {
		log.Fatalf("cluster %s is already running", cluster.Name)
	}

	port, err := util.FindFreePort()
	if err != nil {
		log.WithField("err", err).Fatal("failed to find a free port")
	}

	log.WithFields(log.Fields{
		"pgdata": cluster.Path(),
		"port":   port,
	}).Info("start a postgresql process")
	if err := cluster.Start(port); err != nil {
		log.WithField("err", err).Fatal("failed to start a postgresql process")
	}
}

func ClusterStartCompletion(c *cli.Context) {
	for _, c := range AllClusters() {
		fmt.Println(c.Name)
	}
}
