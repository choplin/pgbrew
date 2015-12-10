package main

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

// DoClusterEnv is an implementation of cluster env command
func DoClusterEnv(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		showHelpAndExit(c, "<cluster> must be specified")
	}

	cluster, err := NewCluster(args[0])
	if err != nil {
		log.WithField("err", err).Fatal("failed to get a cluster")
	}

	// exluce a PATH entry which hash been already set using pgenv
	format := `export PATH=%s:$(echo $PATH | sed -e 's,%s/[^/]*/bin:,,')
export PGDATA=%s
`
	path := filepath.Join(cluster.Pg.Version().Path(), "bin")
	fmt.Printf(format, path, baseDir.installDir(), cluster.Path())

	if cluster.IsRunning() {
		fmt.Printf("export PGPORT=%d\n", cluster.Port)
	}
}

// ClusterEnvCompletion provides cli completion of cluster env command
func ClusterEnvCompletion(c *cli.Context) {
	if len(c.Args()) == 0 {
		for _, c := range AllClusters() {
			fmt.Println(c.Name)
		}
	}
}
