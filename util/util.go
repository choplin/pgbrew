package util

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func RunCommandWithDebugLog(cmd *exec.Cmd) error {
	var waitStdout chan bool
	var waitStderr chan bool

	isDebug := log.GetLevel() == log.DebugLevel

	if isDebug {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}

		waitStdout = handleOutput(stdout, func(line string) {
			log.Debug(line)
		})
		waitStderr = handleOutput(stderr, func(line string) {
			log.Debug(line)
		})
	}

	log.Debug(strings.Join(cmd.Args, " "))
	if err := cmd.Start(); err != nil {
		return err
	}

	if isDebug {
		<-waitStdout
		<-waitStderr
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func handleOutput(out io.ReadCloser, f func(string)) chan bool {
	wait := make(chan bool)
	go func() {
		defer close(wait)

		scanner := bufio.NewScanner(out)

		for scanner.Scan() {
			f(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			//TODO
		}
	}()
	return wait
}
