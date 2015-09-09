package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	home "github.com/mitchellh/go-homedir"
)

const (
	appName    = "pgenv"
	appUsage   = "Manage PostgreSQL environment"
	appVersion = "0.1-SNAPSHOT"
	appAuthor  = "Akihiro Okuno"
	appEmail   = "choplin.choplin@gmail.com"
)

const configFilePathEnv = "PGENV_CONFIG"

const configFilePathSuffix = ".pgenv/config.json"

var (
	configFilePath string
	config         *Config
)

var (
	baseDirectory   string
	localRepository string
	installBase     string
	clusterBase     string
)

func init() {
	configFilePath = os.Getenv(configFilePathEnv)
	if configFilePath == "" {
		configFilePath = filepath.Join(getHomeDir(), configFilePathSuffix)
	}
	config = getConfig(configFilePath)
	if config != nil {
		baseDirectory = config.BasePath
		installBase = filepath.Join(baseDirectory, "versions")
		clusterBase = filepath.Join(baseDirectory, "clusters")
		if config.RepositoryPath == "" {
			localRepository = getLocalRepositoryPath(baseDirectory)
		} else {
			localRepository = config.RepositoryPath
		}
	}
}

func main() {
	app := makeApp()
	app.Run(os.Args)
}

func makeApp() *cli.App {
	app := cli.NewApp()
	app.Author = appAuthor
	app.Email = appEmail
	app.Name = appName
	app.Usage = appUsage
	app.Version = appVersion
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "Enable debug output",
		},
	}

	app.Before = func(c *cli.Context) error {
		args := c.Args()
		if len(args) > 0 && args[0] != "init" {
			if config == nil {
				log.Fatal("pgenv is not initilized. Run `pgenv init` first.")
			}
		}

		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		if len(args) > 0 {
			updateCommandHelp(args[0], commandHelps)
		}

		return nil
	}

	app.Commands = commands

	return app
}

func makeTestEnv() *cli.App {
	app := makeApp()

	log.SetLevel(log.FatalLevel)

	return app
}

func getHomeDir() string {
	home, err := home.Dir()
	if err != nil {
		log.WithField("err", err).Fatal("failed to deterine a home directory")
	}
	return home
}

func getConfig(path string) *Config {
	// This may occur in init command
	if !exists(path) {
		return nil
	}

	config, err := ReadConfigFile(path)
	if err != nil {
		log.WithField("err", err).Fatal("failed to read a config file")
	}
	return config
}

func getLocalRepositoryPath(baseDirectory string) string {
	return filepath.Join(baseDirectory, "repository")
}

func showHelpAndExit(c *cli.Context, msg string) {
	fmt.Fprintf(os.Stderr, "%s\n\n", msg)
	cli.ShowCommandHelp(c, c.Command.Name)
	os.Exit(1)
}

func exists(filename string) bool {
	_, err := os.Lstat(filename)
	return err == nil
}

func updateCommandHelp(command string, helps map[string]string) {
	if help, ok := helps[command]; ok {
		cli.CommandHelpTemplate = `NAME:
   {{.FullName}} - {{.Usage}}

USAGE:
   pgenv {{.FullName}} ` + help + `{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .Flags}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`
	}
}
