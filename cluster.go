package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const clusterExtraInfoFile = ".pgenv_info"
const clusterPortFile = ".pgenv_port"

const defaultPort = 5432

var portConfigRegexp = regexp.MustCompile(`^\s*port\s*=\s*(\d+)`)

type Cluster struct {
	Name string
	Port int
	Pg   *Postgres
}

func NewCluster(name string) (*Cluster, error) {
	c := &Cluster{Name: name}
	if !exists(c.Path()) {
		return nil, fmt.Errorf("cluster %s does not exist", name)
	}

	if err := c.readExtraInfoFile(); err != nil {
		return nil, err
	}
	if err := c.readPortFile(); err != nil {
		return nil, err
	}
	return c, nil
}

func AllClusters() []*Cluster {
	fis, err := ioutil.ReadDir(clusterBase)
	if err != nil {
		log.WithField("err", err).Fatal("failed to get all clusters")
	}

	clusters := make([]*Cluster, len(fis))
	for i, fi := range fis {
		c, err := NewCluster(fi.Name())
		if err != nil {
			log.WithField("err", err).Fatal("failed to get all clusters")
		}
		clusters[i] = c
	}
	return clusters
}

func (c *Cluster) Start(port int) error {
	if err := c.Pg.Start(c.Path(), port); err != nil {
		return err
	}
	c.Port = port
	if err := c.writePortFile(); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) Stop() error {
	if err := c.Pg.Stop(c.Path()); err != nil {
		return err
	}
	c.Port = 0
	return nil
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

func (c *Cluster) IsRunning() bool {
	pidFile := filepath.Join(c.Path(), "postmaster.pid")
	return exists(pidFile)
}

func (c *Cluster) Pid() (int, error) {
	if !c.IsRunning() {
		return 0, nil
	}
	pidFile := filepath.Join(c.Path(), "postmaster.pid")
	file, err := os.Open(pidFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(line))
}

func (c *Cluster) readExtraInfoFile() error {
	path := c.ExtraInfoFilePath()
	out, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	Pg, err := NewPostgres(string(out))
	if err != nil {
		return err
	}

	c.Pg = Pg
	return nil
}

func (c *Cluster) readPortFromPostgresqlConf() (int, error) {
	path := filepath.Join(c.Path(), "postgresql.conf")
	out, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	if portConfigRegexp.Match(out) {
		match := portConfigRegexp.FindSubmatch(out)
		port, err := strconv.Atoi(string(match[1]))
		if err != nil {
			return 0, err
		}
		return port, nil
	}
	return defaultPort, nil
}

func (c *Cluster) portFilePath() string {
	return filepath.Join(c.Path(), clusterPortFile)
}

func (c *Cluster) writePortFile() error {
	path := c.portFilePath()
	str := strconv.Itoa(c.Port)
	return ioutil.WriteFile(path, []byte(str), 0600)
}

func (c *Cluster) readPortFile() error {
	if !c.IsRunning() {
		return nil
	}

	path := c.portFilePath()
	out, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	c.Port, err = strconv.Atoi(string(out))
	if err != nil {
		return err
	}
	return nil
}
