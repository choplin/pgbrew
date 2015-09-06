package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/choplin/pgbrew/git"
	"github.com/choplin/pgbrew/util"
	"github.com/codegangsta/cli"
)

func DoInstall(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
		showHelpAndExit(c, "<version> must be specified")
	}
	version := args[0]

	debug := c.Bool("debug")
	name := c.String("name")
	if name == "" {
		name = defaultName(version, debug)
	}
	options := c.String("options")
	installPath := filepath.Join(installBase, name)
	parallel := c.Bool("parallel")

	hash := doCheckout(version)

	configure(installPath, options, debug)

	makeClean()

	makeInstall(parallel)

	writeVersionFile(name, version, hash)
}

func InstallCompletion(c *cli.Context) {
	repo, err := git.NewRepository(localRepository)
	if err != nil {
		log.WithField("err", err).Fatal("failed to initialize local reporitory")
	}
	tags, err := repo.Tags()
	if err != nil {
		log.WithField("err", err).Fatal("failed to get tags")
	}
	for _, t := range tags {
		fmt.Println(t)
	}
}

func doCheckout(version string) string {
	log.WithField("version", version).Info("git checkout")

	repo, err := git.NewRepository(localRepository)
	if err != nil {
		log.WithField("err", err).Fatal("failed to initialize local reporitory")
	}

	if out, err := repo.Checkout(version); err != nil {
		log.WithFields(log.Fields{
			"version": version,
			"err":     err.Error(),
			"out":     out,
		}).Fatal("failed to checkout")
	}
	hash, err := repo.HeadHash()

	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
			"out": hash,
		}).Fatal("failed to get a hash of HEAD")
	}

	return hash
}

func configure(installPath string, options string, debug bool) {
	cmd := configureCommand(installPath, options, debug)

	log.WithField("options", cmd.Args[1:]).Info("configure")

	err := util.RunCommandWithDebugLog(cmd)
	if err != nil {
		log.WithFields(log.Fields{
			"options": cmd.Args[1:],
			"err":     err.Error(),
		}).Fatal("failed to configure")
	}
}

func makeInstall(parallel bool) {
	args := []string{"install"}
	if parallel {
		args = append(args, "-j")
		args = append(args, fmt.Sprint(runtime.NumCPU()))
	}
	cmd := exec.Command("make", args...)
	cmd.Dir = localRepository

	log.WithField("options", cmd.Args[2:]).Info("make install")
	if err := util.RunCommandWithDebugLog(cmd); err != nil {
		log.WithField("err", err).Fatal("failed to make install")
	}
}

func makeClean() {
	log.Info("make clean")
	cmd := exec.Command("make", "clean")
	cmd.Dir = localRepository
	if err := util.RunCommandWithDebugLog(cmd); err != nil {
		log.WithField("err", err).Fatal("failed to make clean")
	}
}

func writeVersionFile(name string, version string, hash string) {
	installedVersion := Version{
		Name:    name,
		Version: version,
		Hash:    hash,
	}

	log.WithFields(log.Fields{
		"path":    installedVersion.VersionFilePath(),
		"version": version,
		"hash":    hash,
	}).Info("write version file")

	if err := installedVersion.WriteVersionFile(); err != nil {
		log.WithField("err", err).Fatal("failed to write version file")
	}
}

func configureCommand(path string, options string, debug bool) *exec.Cmd {
	args := []string{"--prefix", path}
	if options != "" {
		for _, o := range strings.Split(options, " ") {
			args = append(args, o)
		}
	}
	if debug {
		args = append(args, "--enable-debug")
		args = append(args, "--enable-cassert")
	}

	cmd := exec.Command("./configure", args...)
	cmd.Dir = localRepository

	return cmd
}

func defaultName(version string, debug bool) string {
	var name string
	if strings.HasPrefix(version, "REL") {
		s := strings.Split(version[3:], "_")
		name = fmt.Sprintf("%s.%s.%s", s[0], s[1], s[2])
	} else {
		name = version
	}
	if debug {
		name += "-debug"
	}
	return name
}
