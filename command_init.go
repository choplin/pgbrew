package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func DoInit(c *cli.Context) {
	if exists(configFilePath) {
		fmt.Println("already initilized")
		os.Exit(0)
	}

	log.Info("initialize pgenv")

	basePath := c.String("base-path")
	if basePath == "" {
		basePath = filepath.Join(getHomeDir(), ".pgenv")
	}
	if exists(basePath) {
		log.Fatalf("base-path %s already exists", basePath)
	}

	dirs := []string{
		filepath.Dir(configFilePath),
		basePath,
	}

	for _, d := range dirs {
		if !exists(d) {
			log.WithField("path", d).Debug("create a directory")
			if err := os.Mkdir(d, 0755); err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"path": d,
				}).Fatal("failed to make a directory")
			}
		}
	}

	config := &Config{
		BasePath: basePath,
	}

	log.WithFields(log.Fields{
		"config": config,
		"path":   configFilePath,
	}).Debug("write config file")

	if err := config.Write(configFilePath); err != nil {
		log.WithField("err", err).Fatal("failed to write a config file")
	}
}
