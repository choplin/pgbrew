package main

import "path/filepath"

// BaseDir represents a base directory of pgenv environment
type BaseDir string

func (b BaseDir) installDir() string {
	return filepath.Join(string(b), "versions")
}

func (b BaseDir) clusterDir() string {
	return filepath.Join(string(b), "clusters")
}

func (b BaseDir) defaultLocalRepository() string {
	return filepath.Join(string(b), "repository")
}
