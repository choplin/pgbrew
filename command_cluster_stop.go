package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

// DoClusterStop is an implementation of cluster stop command
func DoClusterStop(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		showHelpAndExit(c, "<cluster> must be specified")
	}

	cluster, err := NewCluster(args[0])
	if err != nil {
		log.WithField("err", err).Fatal("failed to get a cluster")
	}

	if !cluster.IsRunning() {
		log.Fatalf("cluster %s is not running", cluster.Name)
	}

	log.WithFields(log.Fields{
		"pgdata": cluster.Path(),
	}).Debug("stop a postgresql process")
	if err := cluster.Stop(); err != nil {
		log.WithField("err", err).Fatal("failed to start a postgresql process")
	}
}

// ClusterStopCompletion provides cli completion of cluster stop command
func ClusterStopCompletion(c *cli.Context) {
	if len(c.Args()) == 0 {
		for _, c := range AllClusters() {
			fmt.Println(c.Name)
		}
	}
}
