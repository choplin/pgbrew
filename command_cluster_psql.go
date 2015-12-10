package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

// DoClusterPsql is an implementation of cluster psql command
func DoClusterPsql(c *cli.Context) {
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

	var psqlArgs []string
	if len(args) > 1 {
		psqlArgs = args[1:]
	} else {
		psqlArgs = []string{}
	}

	cluster.Pg.Psql(cluster.Port, psqlArgs)
}

// ClusterPsqlCompletion provides cli completion of cluster psql command
func ClusterPsqlCompletion(c *cli.Context) {
	if len(c.Args()) == 0 {
		for _, c := range AllClusters() {
			fmt.Println(c.Name)
		}
	}
}
