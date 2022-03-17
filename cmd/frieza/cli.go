package main

import (
	"fmt"
	"io/ioutil"
	"log"

	. "github.com/outscale-dev/frieza/internal/common"
	"github.com/teris-io/cli"
)

type CliOptions struct {
	debug bool
}

var GlobalCliOptions = CliOptions{
	debug: false,
}

func setupDebug(options map[string]string) {
	if _, ok := options["debug"]; ok {
		GlobalCliOptions.debug = true
	}
}

func cliConfigPath() cli.Option {
	return cli.NewOption("config", "path to configuration file")
}

func cliDebug() cli.Option {
	return cli.NewOption("debug", "enable verbose output for debuging purpose").WithType(cli.TypeBool)
}

func cliJson() cli.Option {
	return cli.NewOption("json", "output in json format (with --plan option only)").WithType(cli.TypeBool)
}

func cliFatalf(json bool, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if json {
		log.Fatalf("{\"error\": \"%s\"}", msg)
	} else {
		log.Fatalf(msg)
	}
}

func cliRoot() cli.App {
	return cli.New(ShortDescription()).
		WithCommand(cliProfile()).
		WithCommand(cliSnapshot()).
		WithCommand(cliClean()).
		WithCommand(cliNuke()).
		WithCommand(cliProvider()).
		WithCommand(cliConfig()).
		WithCommand(cliVersion())
}

func cliVersion() cli.Command {
	return cli.NewCommand("version", "show version").
		WithAction(func(args []string, options map[string]string) int {
			log.Println(FullVersion())
			return 0
		})
}

func ShortDescription() string {
	return "Cleanup your cloud ressources.\n\n" +
		"    Frieza can remove all resources from a cloud account or resources which are not part of a \"snapshot\".\n" +
		"    Snapshots are only a listing of cloud resources.\n" +
		"    Start by adding a new cloud profile with `profile new` sub-command.\n"
}

func disableLogs() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}
