package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

// DoUninstall is an implementation of uninstall command
func DoUninstall(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
		showHelpAndExit(c, fmt.Sprint("<version name> must be specified"))
	}

	name := args[0]
	version, err := NewVersion(name)
	if err != nil {
		log.WithField("err", err).Fatal("failed to get version")
	}

	log.WithField("name", name).Info("uninstall")
	if err := os.RemoveAll(version.Path()); err != nil {
		log.WithField("err", err).Fatal("failed to remove version")
	}
}

// UninstallCompletion provides cli completion of uninstall command
func UninstallCompletion(c *cli.Context) {
	if len(c.Args()) == 0 {
		versions := AllVersions()
		for _, v := range versions {
			fmt.Println(v.Name)
		}
	}
}
