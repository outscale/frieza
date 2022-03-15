package main

import (
	"fmt"
	"log"

	. "github.com/outscale-dev/frieza/internal/common"
	"github.com/teris-io/cli"
)

func cliConfig() cli.Command {
	return cli.NewCommand("config", "configure frieza options").
		WithCommand(cliConfigLs())
}

func cliConfigLs() cli.Command {
	return cli.NewCommand("list", "list configuration options").
		WithOption(cliConfigPath()).
		WithShortcut("ls").
		WithAction(func(args []string, options map[string]string) int {
			configLs(options["config"])
			return 0
		})
}

func configLs(customConfigPath string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoad(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}
	fmt.Println("configuration path:", *configPath)
	fmt.Println("version:", config.Version)
	if len(config.SnapshotFolderPath) == 0 {
		fmt.Println("snapshot_folder_path: (unset)")
	} else {
		fmt.Println("snapshot_folder_path:", config.SnapshotFolderPath)
	}
}
