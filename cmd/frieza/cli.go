package main

import (
	"fmt"

	"github.com/teris-io/cli"
)

var Debug = false

func setupDebug(options map[string]string) {
	if _, ok := options["debug"]; ok {
		Debug = true
	}
}

func cliConfigPath() cli.Option {
	return cli.NewOption("config", "path to configuration file")
}

func cliDebug() cli.Option {
	return cli.NewOption("debug", "enable verbose output for debuging purpose").WithType(cli.TypeBool)
}

func cliRoot() cli.App {
	return cli.New(ShortDescription()).
		WithCommand(cliProfile()).
		WithCommand(cliSnapshot()).
		WithCommand(cliDestroy()).
		WithCommand(cliDestroyAll()).
		WithCommand(cliProvider()).
		WithCommand(cliVersion())
}

func cliVersion() cli.Command {
	return cli.NewCommand("version", "show version").
		WithAction(func(args []string, options map[string]string) int {
			fmt.Println(fullVersion())
			return 0
		})
}

func ShortDescription() string {
	return "Cleanup your cloud ressources.\n\n" +
		"    Frieza can remove all resources from a cloud account or resources which are not part of a \"snapshot\".\n" +
		"    Snapshots are only a listing of cloud resources.\n" +
		"    Start by adding a new cloud profile with `profile new` sub-command.\n"

}
