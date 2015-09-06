package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
)

type JsonListEntry struct {
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	GitRef           string   `json:"git-ref,omitempty"`
	Hash             string   `json:"hash,omitempty"`
	Path             string   `json:"path,omitempty"`
	ConfigureOptions []string `json:"configureOptions,omitempty"`
}

func DoList(c *cli.Context) {
	versions := AllVersions()

	format := c.String("format")
	detail := c.Bool("detail")

	header := []string{"Name", "Version"}
	if detail {
		header = append(header, "Git Reference")
		header = append(header, "Hash")
		header = append(header, "Path")
		header = append(header, "Configure Options")
	}

	switch format {
	case "pretty", "":
		prettyList(header, versions, detail)
	case "plain":
		plainList(header, versions, detail)
	case "json":
		jsonList(versions, detail)
	default:
		showHelpAndExit(c, fmt.Sprint("invalid output format: ", format))
	}
}

func prettyList(header []string, versions []*Version, detail bool) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, v := range versions {
		row := buildRow(v, detail)
		table.Append(row)
	}
	table.Render()
}

func plainList(header []string, versions []*Version, detail bool) {
	fmt.Println(strings.Join(header, "\t"))
	for _, v := range versions {
		row := buildRow(v, detail)
		fmt.Println(strings.Join(row, "\t"))
	}
}

func jsonList(versions []*Version, detail bool) {
	out := []JsonListEntry{}
	for _, v := range versions {
		row := buildRow(v, detail)
		entry := JsonListEntry{
			Name:    row[0],
			Version: row[1],
		}
		if detail {
			entry.GitRef = row[2]
			entry.Hash = row[3]
			entry.Path = row[4]
			entry.ConfigureOptions = strings.Split(row[5], " ")
		}
		out = append(out, entry)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(out)
}

func buildRow(v *Version, detail bool) []string {
	version, err := v.PgVersion()
	if err != nil {
		log.WithField("err", err).Fatal("failed get postgresql version")
	}
	row := []string{v.Name, version}
	if detail {
		d, err := v.Detail()
		if err != nil {
			log.WithFields(log.Fields{
				"version": v.Name,
				"err":     err.Error(),
			}).Fatal("failed to get detailed information")
		}
		row = append(row, v.GitRef)
		row = append(row, v.Hash)
		row = append(row, d.Path)
		row = append(row, strings.Join(d.ConfigureOptions, " "))
	}
	return row
}
