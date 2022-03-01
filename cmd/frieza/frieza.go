package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	app := cliRoot()
	os.Exit(app.Run(os.Args, os.Stdout))
}
