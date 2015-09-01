package util

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
)

type TestFormatter struct{}

func (f *TestFormatter) Format(e *log.Entry) ([]byte, error) {
	return []byte(e.Message + "\n"), nil
}

func TestRunCommandWithDebugLog(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(new(TestFormatter))

	cmd := exec.Command("echo", "foo")
	RunCommandWithDebugLog(cmd)

	b, _ := ioutil.ReadAll(buf)
	got := strings.TrimSpace(string(b))
	want := `echo foo
foo`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
