package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfig_write(t *testing.T) {
	c := &Config{
		BasePath:       "test",
		RepositoryPath: "test-repository",
	}

	dir, _ := ioutil.TempDir("", "config_test")
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "test.config")
	if err := c.Write(path); err != nil {
		t.Fatalf("failed to write a config file: %s", err)
	}

	byte, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read the written config file: %s", err)
	}

	str := string(byte)
	expected := `{
  "base-path": "test",
  "repository-path": "test-repository"
}`
	if str != expected {
		t.Errorf("a wrong output of a config file. got %s, want %s", str, expected)
	}
}

func TestReadConfigFile(t *testing.T) {
	c := &Config{
		BasePath: "test",
	}

	dir, _ := ioutil.TempDir("", "config_test")
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "test.config")
	if err := c.Write(path); err != nil {
		t.Fatalf("failed to write a config file: %s", err)
	}

	got, err := ReadConfigFile(path)
	if err != nil {
		t.Fatalf("failed read a config file")
	}

	if !reflect.DeepEqual(c, got) {
		t.Errorf("readConfigFile returns wrong result. got %s, want %s", got, c)
	}
}
