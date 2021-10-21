package main

import (
	"fmt"
	"log"
	"os"
)

var version string
var commit string

func fullVersion() string {
	if len(version) == 0 {
		version = "0.0.0-beta-unknown-version"
	}
	if len(commit) == 0 {
		commit = "unknown git commit"
	}
	return fmt.Sprintf("%s (%s)", version, commit)
}

func main() {
	log.SetFlags(0)
	app := cliRoot()
	os.Exit(app.Run(os.Args, os.Stdout))
}
