package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func DoClusterEdit(c *cli.Context) {
	args := c.Args()
	if len(args) != 2 {
		showHelpAndExit(c, "<cluster> and <file> must be specified")
	}

	cluster, err := NewCluster(args[0])
	if err != nil {
		log.WithField("err", err).Fatal("failed to get a cluster")
	}

	path := filepath.Join(cluster.Path(), args[1])

	editor := os.Getenv("EDITOR")
	if editor == "" {
		log.Fatal("$EDITOR is not set")
	}

	cmd := exec.Command(editor, path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
}

func ClusterEditCompletion(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		for _, c := range AllClusters() {
			fmt.Println(c.Name)
		}
	} else if len(args) == 1 {
		cluster, err := NewCluster(args[0])
		if err != nil {
			os.Exit(1)
		}
		list, err := listDirectory(cluster.Path())
		for _, e := range list {
			fmt.Println(e)
		}
	}
}

func listDirectory(dir string) ([]string, error) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	ret := make([]string, len(fis))
	for i, fi := range fis {
		ret[i] = fi.Name()
		if fi.IsDir() {
			ret[i] = ret[i] + "/"
		}
	}
	return ret, nil
}
