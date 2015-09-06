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
	Hash             string   `json:"hash,omitempty"`
	Path             string   `json:"path,omitempty"`
	ConfigureOptions []string `json:"configureOptions,omitempty"`
}

func DoList(c *cli.Context) {
	versions := AllVersions()

	format := c.String("format")
	detail := c.Bool("detail")

	switch format {
	case "pretty", "":
		prettyList(versions, detail)
	case "plain":
		plainList(versions, detail)
	case "json":
		jsonList(versions, detail)
	default:
		showHelpAndExit(c, fmt.Sprint("invalid output format: ", format))
	}
}

func prettyList(versions []*Version, detail bool) {
	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"Name", "Version"}
	if detail {
		header = append(header, "Hash")
		header = append(header, "Path")
		header = append(header, "Configure Options")
	}
	table.SetHeader(header)

	for _, v := range versions {
		row := []string{v.Name, v.Version}
		if detail {
			d, err := v.Detail()
			if err != nil {
				log.WithFields(log.Fields{
					"version": v.Name,
					"err":     err.Error(),
				}).Fatal("failed to get detailed information")
			}
			row = append(row, v.Hash)
			row = append(row, d.Path)
			row = append(row, strings.Join(d.ConfigureOptions, " "))
		}
		table.Append(row)
	}
	table.Render()
}

func plainList(versions []*Version, detail bool) {
	for _, v := range versions {
		row := []string{v.Name, v.Version}
		if detail {
			d, err := v.Detail()
			if err != nil {
				log.WithFields(log.Fields{
					"version": v.Name,
					"err":     err.Error(),
				}).Fatal("failed to get detailed information")
			}
			row = append(row, v.Hash)
			row = append(row, d.Path)
			row = append(row, strings.Join(d.ConfigureOptions, " "))
		}
		fmt.Println(strings.Join(row, "\t"))
	}
}

func jsonList(versions []*Version, detail bool) {
	out := []JsonListEntry{}
	for _, v := range versions {
		entry := JsonListEntry{
			Name:    v.Name,
			Version: v.Version,
		}
		if detail {
			d, err := v.Detail()
			if err != nil {
				log.WithFields(log.Fields{
					"version": v.Name,
					"err":     err.Error(),
				}).Fatal("failed to get detailed information")
			}
			entry.Hash = v.Hash
			entry.Path = d.Path
			entry.ConfigureOptions = d.ConfigureOptions
		}
		out = append(out, entry)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(out)
}
