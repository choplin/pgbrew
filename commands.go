package main

import "github.com/codegangsta/cli"

var tasks = []string{"cook", "clean", "laundry", "eat", "sleep", "code"}

var commandHelps = map[string]string{
	"init":      "[-p PATH]",
	"clone":     "[-o OPTIONS]",
	"available": "",
	"install":   "[-n NAME] [-d] [-o OPTOINS] [-p] <tag|branch|commit>",
	"list":      "[-f pretty|plain|json] [-d]",
	"uninstall": "<version>",
	"current":   "[-u] | <version>",
	"initdb":    "<version> <name> [-o OPTIONS]",
}

var commands = []cli.Command{
	initCommand,
	cloneCommand,
	updateCommand,
	availableCommand,
	installCommand,
	listCommand,
	uninstallCommand,
	currentCommand,
	initdbCommand,
	clusterCommand,
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Initialize pgbrew environment",
	Description: `During initialization process, a config file will be created at ~/.pgbrew/config.json.
   This path is not related to a pgbrew base directory, and not configurable.`,
	Action: DoInit,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path,p",
			Usage: "Path of pgbrew base directory. default: ~/.pgbrew",
		},
	},
}

var cloneCommand = cli.Command{
	Name:   "clone",
	Usage:  "Clone PostgreSQL git repository into a local directory",
	Action: DoClone,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "options,o",
			Usage: "Options passed to git clone",
		},
	},
}

var updateCommand = cli.Command{
	Name:   "update",
	Usage:  "Update a local git repository",
	Action: DoUpdate,
}

var availableCommand = cli.Command{
	Name:   "available",
	Usage:  "List available versions",
	Action: DoAvailable,
}

var installCommand = cli.Command{
	Name:  "install",
	Usage: "Build and install a specified version",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name,n",
			Usage: "A name for this installation. a default value is a version number in x.y.z format.",
		},
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "Enable debug build swith (i.e. --enable-debug --enable-cassert)",
		},
		cli.StringFlag{
			Name:  "options,o",
			Usage: "Options passed to configure",
		},
		cli.BoolFlag{
			Name:  "parallel,p",
			Usage: "Allow multiple jobs of make command",
		},
	},
	Action:       DoInstall,
	BashComplete: InstallCompletion,
}

var listCommand = cli.Command{
	Name:  "list",
	Usage: "List installed versions",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format,f",
			Usage: "Output format. available option: pretty plain json. default: pretty.",
		},
		cli.BoolFlag{
			Name:  "detail,d",
			Usage: "Extend output with detailed information",
		},
	},
	Action: DoList,
}

var uninstallCommand = cli.Command{
	Name:         "uninstall",
	Usage:        "Uninstall a specified version",
	Action:       DoUninstall,
	BashComplete: UninstallCompletion,
}

var currentCommand = cli.Command{
	Name:  "current",
	Usage: "Set or show the current version",
	Description: `If <version> is not specified, display a current version.
   If -u option is specified, unset a current version.
   Otherwise, set a specfied version as a current version.`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "unset,u",
			Usage: "Unset a current version",
		},
	},
	Action:       DoCurrent,
	BashComplete: CurrentCompletion,
}

var initdbCommand = cli.Command{
	Name:  "initdb",
	Usage: "Execute initdb to create a new cluster",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name,n",
			Usage: "A name for this cluster. A specified <version> is used as a default value.",
		},
		cli.StringFlag{
			Name:  "options,o",
			Usage: "Options passed to configure",
		},
	},
	Action:       DoInitdb,
	BashComplete: InitdbCompletion,
}

var clusterCommands = []cli.Command{
	clusterRemoveCommand,
	clusterListCommand,
	clusterEnvCommand,
	clusterStartCommand,
	clusterStopCommand,
	clusterPSQLCommand,
}

var clusterCommand = cli.Command{
	Name:        "cluster",
	Usage:       "Manage PostgreSQL clusters",
	Subcommands: clusterCommands,
}

var clusterRemoveCommand = cli.Command{
	Name:         "remove",
	Usage:        "Remove a specified cluster",
	Action:       DoClusterRemove,
	BashComplete: ClusterRemoveCompletion,
}

var clusterListCommand = cli.Command{
	Name:  "list",
	Usage: "List clusters",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format,f",
			Usage: "Output format. available option: pretty plain json. default: pretty.",
		},
		cli.BoolFlag{
			Name:  "detail,d",
			Usage: "Extend output with detailed information",
		},
	},
	Action: DoClusterList,
}

var clusterEnvCommand = cli.Command{
	Name:   "env",
	Usage:  "Show shell scripts to enable a specified cluster",
	Action: DoClusterEnv,
}

var clusterStartCommand = cli.Command{
	Name:   "start",
	Usage:  "Start a postgresql process with a specified cluster",
	Action: DoClusterStart,
}

var clusterStopCommand = cli.Command{
	Name:   "stop",
	Usage:  "Stop a postgresql process with a specified cluster",
	Action: DoClusterStop,
}

var clusterPSQLCommand = cli.Command{
	Name:   "psql",
	Usage:  "Run psql for a specified cluster",
	Action: DoClusterPSQL,
}
