package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/choplin/pgenv/git"
	"github.com/choplin/pgenv/util"
	"github.com/codegangsta/cli"
)

func DoInstall(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		showHelpAndExit(c, "<tag|branch|commit> must be specified")
	}
	gitRef := args[0]
	var configureOptions []string
	if len(args) > 1 {
		configureOptions = args[1:]
	} else {
		configureOptions = []string{}
	}

	debug := c.Bool("debug")
	name := c.String("name")
	if name == "" {
		name = defaultName(gitRef, debug)
	}
	installPath := filepath.Join(installBase, name)
	parallel := c.Bool("parallel")

	hash := doCheckout(gitRef)

	configure(installPath, configureOptions, debug)

	makeClean()

	makeInstall(parallel)

	WriteExtraInfoFile(name, gitRef, hash)
}

func InstallCompletion(c *cli.Context) {
	if len(c.Args()) == 0 {
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
}

func doCheckout(gitRef string) string {
	log.WithField("git ref", gitRef).Info("git checkout")

	repo, err := git.NewRepository(localRepository)
	if err != nil {
		log.WithField("err", err).Fatal("failed to initialize local reporitory")
	}

	if out, err := repo.Checkout(gitRef); err != nil {
		log.WithFields(log.Fields{
			"git ref": gitRef,
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

func configure(installPath string, options []string, debug bool) {
	cmd := configureCommand(installPath, options, debug)

	log.WithField("configure options", cmd.Args[1:]).Info("configure")

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

func WriteExtraInfoFile(name string, gitRef string, hash string) {
	installedVersion := Version{
		Name:   name,
		GitRef: gitRef,
		Hash:   hash,
	}

	log.WithFields(log.Fields{
		"path":    installedVersion.ExtraInfoFilePath(),
		"git ref": gitRef,
		"hash":    hash,
	}).Info("write an extra info file")

	if err := installedVersion.WriteExtraInfoFile(); err != nil {
		log.WithField("err", err).Fatal("failed to write an extra info file")
	}
}

func configureCommand(path string, options []string, debug bool) *exec.Cmd {
	args := []string{"--prefix", path}
	for _, o := range options {
		args = append(args, o)
	}
	if debug {
		args = append(args, "--enable-debug")
		args = append(args, "--enable-cassert")
	}

	cmd := exec.Command("./configure", args...)
	cmd.Dir = localRepository

	return cmd
}

func defaultName(gitRef string, debug bool) string {
	var name string
	if strings.HasPrefix(gitRef, "REL") {
		s := strings.Split(gitRef[3:], "_")
		name = fmt.Sprintf("%s.%s.%s", s[0], s[1], s[2])
	} else {
		name = gitRef
	}
	if debug {
		name += "-debug"
	}
	return name
}
