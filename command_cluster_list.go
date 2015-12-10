package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
)

// ClusterRow represents a row of output of cluster list command
type ClusterRow struct {
	Name        string `json:"name"`
	VersionName string `json:"version-name"`
	State       string `json:"state"`
	Port        int    `json:"port,omitempty"`
	Pid         int    `json:"pid,omitempty"`
	Path        string `json:"path,omitempty"`
	detail      bool
}

var (
	commonheader = []string{"Name", "Version Name", "State"}
	detailHeader = []string{"Port", "Pid", "Path"}
)

// DoClusterList is an implementation of cluster list command
func DoClusterList(c *cli.Context) {
	clusters := AllClusters()

	format := c.String("format")
	detail := c.Bool("detail")

	header := commonheader
	if detail {
		for _, h := range detailHeader {
			header = append(header, h)
		}
	}

	switch format {
	case "pretty", "":
		prettyClusterList(header, clusters, detail)
	case "plain":
		plainClusterList(header, clusters, detail)
	case "json":
		jsonClusterList(clusters, detail)
	default:
		showHelpAndExit(c, fmt.Sprint("invalid output format: ", format))
	}
}

func prettyClusterList(header []string, clusters []*Cluster, detail bool) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, c := range clusters {
		row := buildClusterRow(c, detail)
		table.Append(row.toStringSlice())
	}
	table.Render()
}

func plainClusterList(header []string, clusters []*Cluster, detail bool) {
	fmt.Println(strings.Join(header, "\t"))
	for _, c := range clusters {
		row := buildClusterRow(c, detail)
		fmt.Println(strings.Join(row.toStringSlice(), "\t"))
	}
}

func jsonClusterList(clusters []*Cluster, detail bool) {
	out := []*ClusterRow{}
	for _, c := range clusters {
		row := buildClusterRow(c, detail)
		out = append(out, row)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(out)
}

func buildClusterRow(c *Cluster, detail bool) *ClusterRow {
	row := &ClusterRow{
		Name:        c.Name,
		VersionName: c.Pg.Version().Name,
		detail:      detail,
	}
	if c.IsRunning() {
		row.State = "running"
	} else {
		row.State = "stopped"
	}

	if detail {
		row.Port = c.Port
		pid, err := c.Pid()
		if err != nil {
			log.WithField("err", err).Fatal("failed to get a pid")
		}
		row.Pid = pid
		row.Path = c.Path()
	}

	return row
}

func (r *ClusterRow) toStringSlice() []string {
	ret := []string{
		r.Name,
		r.VersionName,
		r.State,
	}
	if r.detail {
		if r.Port == 0 {
			ret = append(ret, "")
		} else {
			ret = append(ret, strconv.Itoa(r.Port))
		}

		if r.Pid == 0 {
			ret = append(ret, "")
		} else {
			ret = append(ret, strconv.Itoa(r.Pid))
		}
		ret = append(ret, r.Path)
	}
	return ret
}
