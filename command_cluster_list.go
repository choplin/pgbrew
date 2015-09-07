package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
)

type ClusterRow struct {
	Name        string `json:"name"`
	VersionName string `json:"version-name"`
	State       string `json:"state"`
	detail      bool
}

func DoClusterList(c *cli.Context) {
	clusters := AllClusters()

	format := c.String("format")
	detail := c.Bool("detail")

	header := []string{"Name", "Version Name", "State"}
	if detail {
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
	return row
}

func (r *ClusterRow) toStringSlice() []string {
	ret := []string{
		r.Name,
		r.VersionName,
		r.State,
	}
	if !r.detail {
	}
	return ret
}
