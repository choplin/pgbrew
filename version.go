package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const versionExtraInfoFile = ".pgenv_info"

// Version represents each installed version.
type Version struct {
	Name   string
	GitRef string
	Hash   string
}

// VersionDetail represents detailed information of each installed version. Most of its fields are derived from pg_config.
type VersionDetail struct {
	Path             string
	ConfigureOptions []string
}

// NewVersion instanciate and initialize version with a specified name
func NewVersion(name string) (*Version, error) {
	version := &Version{Name: name}
	path := version.Path()
	if !exists(path) {
		return nil, fmt.Errorf("%s is not installed", name)
	}

	if err := version.readExtraInfoFile(); err != nil {
		return nil, err
	}

	return version, nil
}

// AllVersions lists all installed versions
func AllVersions() []*Version {
	fis, err := ioutil.ReadDir(baseDir.installDir())
	if err != nil {
		log.WithField("err", err).Fatal("failed to get all installed versions")
	}

	versions := make([]*Version, len(fis))
	for i, fi := range fis {
		v, err := NewVersion(fi.Name())
		if err != nil {
			log.WithField("err", err).Fatal("failed to get all installed versions")
		}
		versions[i] = v
	}
	return versions
}

// WriteExtraInfoFile writes extra information of this version into a file
func (v *Version) WriteExtraInfoFile() error {
	path := v.ExtraInfoFilePath()
	str := v.Hash + "\t" + v.GitRef
	return ioutil.WriteFile(path, []byte(str), 0644)
}

// Path returns a path of this version on the filesystem
func (v *Version) Path() string {
	return filepath.Join(baseDir.installDir(), v.Name)
}

// ExtraInfoFilePath returns a path of extra info file on the filesystem
func (v *Version) ExtraInfoFilePath() string {
	return filepath.Join(v.Path(), versionExtraInfoFile)
}

// Detail returns a detailed information of this version
func (v *Version) Detail() (*VersionDetail, error) {
	pg := &Postgres{v}
	configureOut, err := pg.PgConfig("--configure")
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

// PgVersion returns a result of `pg --versions` of this version
func (v *Version) PgVersion() (string, error) {
	pg := &Postgres{v}
	version, err := pg.PgConfig("--version")
	if err != nil {
		return "", err
	}
	return version, nil
}

func (v *Version) readExtraInfoFile() error {
	path := v.ExtraInfoFilePath()
	out, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	info := strings.Split(string(out), "\t")
	v.Hash = info[0]
	v.GitRef = info[1]

	return nil
}
