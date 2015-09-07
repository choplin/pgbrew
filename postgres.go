package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/choplin/pgenv/util"
)

type Postgres struct {
	v *Version
}

func NewPostgres(name string) (*Postgres, error) {
	version, err := NewVersion(name)
	if err != nil {
		return nil, err
	}
	return &Postgres{version}, nil
}

func (p *Postgres) Version() *Version {
	return p.v
}

func (p *Postgres) Initdb(args []string) error {
	bin := p.binPath("initdb")
	cmd := exec.Command(bin, args...)
	return util.RunCommandWithDebugLog(cmd)
}

func (p *Postgres) Start(pgdata string, port int) error {
	bin := p.binPath("pg_ctl")
	args := []string{"start", "-D", pgdata}
	cmd := exec.Command(bin, args...)

	cmd.Env = []string{
		fmt.Sprintf("PGPORT=%d", port),
	}

	return util.RunCommandWithDebugLog(cmd)
}

func (p *Postgres) Stop(pgdata string) error {
	bin := p.binPath("pg_ctl")
	args := []string{"stop", "-D", pgdata}
	cmd := exec.Command(bin, args...)
	return util.RunCommandWithDebugLog(cmd)
}

func (p *Postgres) PgConfig(option string) (string, error) {
	pgConfig := p.binPath("pg_config")
	cmd := exec.Command(pgConfig, option)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func (p *Postgres) binPath(name string) string {
	return filepath.Join(p.v.Path(), "bin", name)
}
