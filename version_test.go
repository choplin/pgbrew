package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewVersion(t *testing.T) {
	_, err := NewVersion("dummy")
	if err == nil {
		t.Errorf("NewVersion with a wrong version must return error")
	}

	version, _ := NewVersion("9.4.4")
	want := &Version{
		Name:   "9.4.4",
		GitRef: "REL9_4_4",
		Hash:   "7c055f3ec3bd338a1ebb8c73cff3d01df626471e",
	}
	if !reflect.DeepEqual(version, want) {
		t.Errorf("NewVersion retured a wrong version. got %s, want %s", version, want)
	}
}

func TestAllVersions(t *testing.T) {
	versions := AllVersions()
	names := []string{"9.3.9-debug", "9.4.4"}

	for i, v := range versions {
		if v.Name != names[i] {
			t.Errorf("AllVersions returned a wrong version. got %s, want %s", v.Name, names[i])
		}
	}
}

func TestVersion_Path(t *testing.T) {
	version, _ := NewVersion("9.4.4")
	path := version.Path()
	want := filepath.Join(baseDir.installDir(), "9.4.4")
	if path != want {
		t.Errorf("Version.Path() returned a wrong path. got %s, want %s", path, want)
	}
}

func TestVersion_VersionFilePath(t *testing.T) {
	version, _ := NewVersion("9.4.4")
	path := version.ExtraInfoFilePath()
	want := filepath.Join(baseDir.installDir(), "9.4.4", versionExtraInfoFile)
	if path != want {
		t.Errorf("Version.ExtraInfoFilePath() returned a wrong path. got %s, want %s", path, want)
	}
}

func TestVersion_Detail(t *testing.T) {
	version, _ := NewVersion("9.4.4")
	got, err := version.Detail()
	if err != nil {
		t.Fatal("failed to get a detail of version")
	}
	want := &VersionDetail{
		Path:             version.Path(),
		ConfigureOptions: []string{"--prefix", "/home/postgres/.pgenv/versions/9.4.4"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Version.Detail() returned a wrong result. got %s, want %s", got, want)
	}
}
