package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const versionFileName = "version"

// Version represents each installed version.
type Version struct {
	Name    string
	Version string
	Hash    string
}

// VersionDetail represents detailed information of each installed version. Most of its fields are derived from pg_config.
type VersionDetail struct {
	Path             string
	ConfigureOptions []string
}

func NewVersion(name string) (*Version, error) {
	version := &Version{Name: name}
	path := version.Path()
	if !exists(path) {
		return nil, fmt.Errorf("%s is not installed", name)
	}

	if err := version.readVersionFile(); err != nil {
		return nil, err
	}

	return version, nil
}

func AllVersions() []Version {
	fis, err := ioutil.ReadDir(installBase)
	if err != nil {
		log.WithField("err", err).Fatal("failed to get all installed versions")
	}

	versions := make([]Version, len(fis))
	for i, fi := range fis {
		v, err := NewVersion(fi.Name())
		if err != nil {
			log.WithField("err", err).Fatal("failed to get all installed versions")
		}
		versions[i] = *v
	}
	return versions
}

func (v *Version) WriteVersionFile() error {
	path := v.VersionFilePath()
	str := v.Hash + "\t" + v.Version
	return ioutil.WriteFile(path, []byte(str), 0644)
}

func (v *Version) Path() string {
	return filepath.Join(installBase, v.Name)
}

func (v *Version) VersionFilePath() string {
	return filepath.Join(v.Path(), versionFileName)
}

func (v *Version) Detail() (*VersionDetail, error) {
	configureOut, err := v.pgConfig("--configure")
	if err != nil {
		return nil, err
	}

	configureOptions := []string{}
	for _, c := range strings.Split(string(configureOut), " ") {
		configureOptions = append(configureOptions, strings.Trim(c, "'"))
	}
	return &VersionDetail{
		Path:             v.Path(),
		ConfigureOptions: configureOptions,
	}, nil
}

func (v *Version) readVersionFile() error {
	path := v.VersionFilePath()
	out, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	info := strings.Split(string(out), "\t")
	v.Hash = info[0]
	v.Version = info[1]

	return nil
}

func (v *Version) pgConfig(option string) (string, error) {
	pgConfig := filepath.Join(v.Path(), "bin", "pg_config")
	cmd := exec.Command(pgConfig, option)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\n"), nil

}
