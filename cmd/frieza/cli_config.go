package main

import (
	"fmt"
	"log"

	. "github.com/outscale-dev/frieza/internal/common"
	"github.com/teris-io/cli"
)

func cliConfig() cli.Command {
	return cli.NewCommand("config", "configure frieza options").
		WithCommand(cliConfigLs()).
		WithCommand(cliConfigSet())
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

func cliConfigSet() cli.Command {
	return cli.NewCommand("set", "set a specific option").
		WithOption(cliConfigPath()).
		WithArg(cli.NewArg("option_name", "option's name to set")).
		WithArg(cli.NewArg("option_value", "option's value to set")).
		WithAction(func(args []string, options map[string]string) int {
			configSet(options["config"], &args[0], &args[1])
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
	fmt.Println("version:", config.Version)
	if len(config.SnapshotFolderPath) == 0 {
		fmt.Println("snapshot_folder_path: (unset)")
	} else {
		fmt.Println("snapshot_folder_path:", config.SnapshotFolderPath)
	}
}

func configSet(customConfigPath string, optionName *string, optionValue *string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatal("Cannot load configuration: " + err.Error())
	}
	switch *optionName {
	case "snapshot_folder_path":
		config.SnapshotFolderPath = *optionValue
	default:
		log.Fatalf("Unknow option name")
	}
	if err = config.Write(configPath); err != nil {
		log.Fatalf("Cannot save configuration file: %s", err.Error())
	}
}
