package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
	installPath := filepath.Join(baseDir.installDir(), name)
	parallel := c.Bool("parallel")

	hash, err := build(gitRef, installPath, configureOptions, debug, parallel)
	if err != nil {
		log.WithField("err", err).Fatal("failed to build")
	}

	WriteExtraInfoFile(name, gitRef, hash)
}

func InstallCompletion(c *cli.Context) {
	if len(c.Args()) == 0 {
		repo, err := git.NewRepository(config.RepositoryPath)
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

func build(gitRef string, installPath string, configureOptions []string, debug bool, parallel bool) (string, error) {
	workdir, err := ioutil.TempDir("", "pgbrew-install")
	if err != nil {
		return "", fmt.Errorf("failed to create a temporary directory: %s", err)
	}
	defer os.RemoveAll(workdir)
	log.WithField("path", workdir).Debug("create a temporary working directory")

	hash, err := doCheckout(gitRef, workdir)
	if err != nil {
		return "", err
	}

	if err := configure(workdir, installPath, configureOptions, debug); err != nil {
		return "", err
	}

	if err := makeInstall(parallel, workdir); err != nil {
		return "", err
	}

	return hash, nil
}

func doCheckout(gitRef string, workdir string) (string, error) {
	log.WithField("git ref", gitRef).Info("git checkout")

	repo, err := git.NewRepository(config.RepositoryPath)
	if err != nil {
		return "", fmt.Errorf("failed to get a repository: %s", err)
	}

	if out, err := repo.CheckoutWithWorkTree(gitRef, workdir); err != nil {
		return "", fmt.Errorf("failed to checkout: %s, %s", err, out)
	}
	defer repo.Checkout("master")

	hash, err := repo.HeadHash()
	if err != nil {
		return "", fmt.Errorf("failed to get a hash of HEAD: %s", err)
	}

	return hash, nil
}

func configure(workdir string, installPath string, options []string, debug bool) error {
	cmd := configureCommand(workdir, installPath, options, debug)

	log.WithField("configure options", cmd.Args[1:]).Info("configure")

	err := util.RunCommandWithDebugLog(cmd)
	if err != nil {
		return err
	}
	return nil
}

func makeInstall(parallel bool, workdir string) error {
	args := []string{"install"}
	if parallel {
		args = append(args, "-j")
		args = append(args, fmt.Sprint(runtime.NumCPU()))
	}
	cmd := exec.Command("make", args...)
	cmd.Dir = workdir

	log.WithField("options", cmd.Args[2:]).Info("make install")
	if err := util.RunCommandWithDebugLog(cmd); err != nil {
		return fmt.Errorf("failed to make install: %s", err)
	}
	return nil
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

func configureCommand(workdir string, path string, options []string, debug bool) *exec.Cmd {
	args := []string{"--prefix", path}
	for _, o := range options {
		args = append(args, o)
	}
	if debug {
		args = append(args, "--enable-debug")
		args = append(args, "--enable-cassert")
	}

	cmd := exec.Command("./configure", args...)
	cmd.Dir = workdir

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
