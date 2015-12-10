package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/choplin/pgenv/util"
	"github.com/codegangsta/cli"
)

// DoClusterStart is an implementation of cluster start command
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

	port := c.Int("port")
	if port == 0 {
		if c.Bool("find-free-port") {
			port, err = util.FindFreePort()
			if err != nil {
				log.WithField("err", err).Fatal("failed to find a free port")
			}
		} else {
			port, err = cluster.readPortFromPostgresqlConf()
			if err != nil {
				log.WithField("err", err).Fatal("failed to read a postgresql conf")
			}
		}
	}

	log.WithFields(log.Fields{
		"pgdata": cluster.Path(),
		"port":   port,
	}).Debug("start a postgresql process")
	if err := cluster.Start(port); err != nil {
		log.WithField("err", err).Fatal("failed to start a postgresql process")
	}
}

// ClusterStartCompletion provides cli completion of cluster start command
func ClusterStartCompletion(c *cli.Context) {
	if len(c.Args()) == 0 {
		for _, c := range AllClusters() {
			fmt.Println(c.Name)
		}
	}
}
