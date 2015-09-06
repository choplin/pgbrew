package main

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/choplin/pgbrew/util"
)

const clusterVersionFileName = ".pgbrew_version"

type Cluster struct {
	pg   *Postgres
	name string
}

func (c *Cluster) Start() error {
	port, err := util.FindFreePort()
	if err != nil {
		log.WithField("err", err).Fatal("failed to find a free port")
	}

	return c.pg.Start(port)
}

func (c *Cluster) WriteVersionFile() error {
	path := c.VersionFilePath()
	str := c.name
	return ioutil.WriteFile(path, []byte(str), 0600)
}

func (c *Cluster) Path() string {
	return filepath.Join(clusterBase, c.name)
}

func (c *Cluster) VersionFilePath() string {
	return filepath.Join(c.Path(), clusterVersionFileName)
}

func (c *Cluster) readVersionFile() error {
	path := c.VersionFilePath()
	out, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	pg, err := NewPostgres(string(out))
	if err != nil {
		return err
	}

	c.pg = pg
	return nil
}
