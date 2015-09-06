package main

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/choplin/pgbrew/util"
)

const clusterExtraInfoFile = ".pgbrew_info"

type Cluster struct {
	Name string
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

func (c *Cluster) WriteExtraInfoFile() error {
	path := c.ExtraInfoFilePath()
	str := c.Name
	return ioutil.WriteFile(path, []byte(str), 0600)
}

func (c *Cluster) Path() string {
	return filepath.Join(clusterBase, c.Name)
}

func (c *Cluster) ExtraInfoFilePath() string {
	return filepath.Join(c.Path(), clusterExtraInfoFile)
}

func (c *Cluster) readExtraInfoFile() error {
	path := c.ExtraInfoFilePath()
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
