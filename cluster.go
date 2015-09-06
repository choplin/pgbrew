package main

import (
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/choplin/pgbrew/util"
)

type Cluster struct {
	version *Version
	pg      *Postgres
	name    string
}

func (c *Cluster) Start() error {
	port, err := util.FindFreePort()
	if err != nil {
		log.WithField("err", err).Fatal("failed to find a free port")
	}

	return c.pg.Start(port)
}

func (c *Cluster) path() string {
	return filepath.Join(clusterBase, c.name)
}
