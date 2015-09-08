package main

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type Config struct {
	BasePath       string `json:"base-path"`
	RepositoryPath string `json:"repository-path,omitempty"`
}

func ReadConfigFile(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Write(path string) error {
	out, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(path, out, 0644); err != nil {
		return err
	}

	return nil
}

func (c *Config) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		log.WithField("err", err).Fatal("failed to convert config to string")
	}
	return string(out)
}
