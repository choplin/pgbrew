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

// Row represents a row of output of list command
type Row struct {
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	GitRef           string   `json:"git-ref,omitempty"`
	Hash             string   `json:"hash,omitempty"`
	Path             string   `json:"path,omitempty"`
	ConfigureOptions []string `json:"configureOptions,omitempty"`
	detail           bool
}

// DoList is an implementation of list command
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
		table.Append(row.toStringSlice())
	}
	table.Render()
}

func plainList(header []string, versions []*Version, detail bool) {
	fmt.Println(strings.Join(header, "\t"))
	for _, v := range versions {
		row := buildRow(v, detail)
		fmt.Println(strings.Join(row.toStringSlice(), "\t"))
	}
}

func jsonList(versions []*Version, detail bool) {
	out := []*Row{}
	for _, v := range versions {
		row := buildRow(v, detail)
		out = append(out, row)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(out)
}

func buildRow(v *Version, detail bool) *Row {
	version, err := v.PgVersion()
	if err != nil {
		log.WithField("err", err).Fatal("failed get postgresql version")
	}
	row := &Row{Name: v.Name, Version: version, detail: detail}
	if detail {
		d, err := v.Detail()
		if err != nil {
			log.WithFields(log.Fields{
				"version": v.Name,
				"err":     err.Error(),
			}).Fatal("failed to get detailed information")
		}
		row.GitRef = v.GitRef
		row.Hash = v.Hash
		row.Path = d.Path
		row.ConfigureOptions = d.ConfigureOptions
	}
	return row
}

func (r *Row) toStringSlice() []string {
	if !r.detail {
		return []string{
			r.Name,
			r.Version,
		}
	}

	return []string{
		r.Name,
		r.Version,
		r.GitRef,
		r.Hash,
		r.Path,
		strings.Join(r.ConfigureOptions, " "),
	}
}
