package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	. "github.com/outscale-dev/frieza/internal/common"
	"github.com/teris-io/cli"
)

func cliSnapshot() cli.Command {
	return cli.NewCommand("snapshot", "manage resource snapshots").
		WithShortcut("snap").
		WithCommand(cliSnapshotNew()).
		WithCommand(cliSnapshotLs()).
		WithCommand(cliSnapshotDescribe()).
		WithCommand(cliSnapshotRm())
}

func cliSnapshotNew() cli.Command {
	return cli.NewCommand("new", "create new snapshot containing all resource ids").
		WithArg(cli.NewArg("snapshot_name", "snapshot name")).
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithArg(cli.NewArg("profile", "one or more profile to snapshot").AsOptional()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			snapshotNew(options["config"], args)
			return 0
		})
}

func cliSnapshotLs() cli.Command {
	return cli.NewCommand("list", "list snapshots").
		WithShortcut("ls").
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			snapshotLs(options["config"])
			return 0
		})
}

func cliSnapshotDescribe() cli.Command {
	return cli.NewCommand("describe", "describe snapshot").
		WithShortcut("desc").
		WithArg(cli.NewArg("snapshot_name", "snapshot name to describe")).
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			snapshotDescribe(options["config"], &args[0])
			return 0
		})
}

func cliSnapshotRm() cli.Command {
	return cli.NewCommand("remove", "remove snapshot").
		WithArg(cli.NewArg("snapshot_name", "snapshot's name to remove")).
		WithShortcut("rm").
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			snapshotRm(options["config"], &args[0])
			return 0
		})
}

func snapshotNew(customConfigPath string, args []string) {
	if len(args) < 2 {
		log.Fatal("No profile has been chosen to be snapshoted")
	}
	snapshotName := args[0]
	profileNames := args[1:]
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}
	if _, err = SnapshotLoad(snapshotName, config); err == nil {
		log.Fatalf("Snapshot %s already exist", snapshotName)
	}

	var providers []Provider
	var profiles []string
	for _, profileName := range profileNames {
		found := false
		for _, profile := range config.Profiles {
			if profileName == profile.Name {
				provider, err := ProviderNew(profile)
				if err != nil {
					log.Fatalf("Cannot initialize profile %s: %s", profile.Name, err.Error())
				}
				providers = append(providers, provider)
				profiles = append(profiles, profile.Name)
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("Profile %s not found", profileName)
		}
	}

	for i, provider := range providers {
		if err := provider.AuthTest(); err != nil {
			log.Fatalf("Provider test failed for profile %s: %s", profiles[i], err.Error())
		}
	}

	date := fmt.Sprintf("%s", time.Now().UTC())
	snapshot := Snapshot{
		Version: SnapshotVersion(),
		Name:    snapshotName,
		Date:    date,
		Config:  config,
	}
	for i, provider := range providers {
		snapshot.Data = append(snapshot.Data, SnapshotData{
			Profile: profiles[i],
			Objects: ReadObjects(&provider),
		})
	}
	if err = snapshot.Write(); err != nil {
		log.Fatalf("Snapshot failed: %s", err.Error())
	}
}

func snapshotLs(customConfigPath string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}
	files, err := ioutil.ReadDir(config.SnapshotFolderPath)
	if err != nil {
		log.Fatalf("Error while listing snapshots: %s", err.Error())
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		snapshotName := strings.TrimSuffix(file.Name(), ".json")
		if snapshot, err := SnapshotLoad(snapshotName, config); err == nil {
			fmt.Println(snapshot.Name)
		}
	}
}

func snapshotDescribe(customConfigPath string, snapshotName *string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}
	if err != nil {
		log.Fatalf("Error while reading snapshots: %s", err.Error())
	}
	snapshot, err := SnapshotLoad(*snapshotName, config)
	if err != nil {
		log.Fatalf("Cannot load snapshot %s: %s", *snapshotName, err.Error())
	}
	fmt.Print(snapshot)
}

func snapshotRm(customConfigPath string, snapshotName *string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}
	snapshot, err := SnapshotLoad(*snapshotName, config)
	if err != nil {
		log.Fatalf("Error while deleting snapshot %s: %s", *snapshotName, err.Error())
	}
	if err = snapshot.Delete(); err != nil {
		log.Fatalf("Error while deleting snapshot %s: %s", *snapshotName, err.Error())
	}
}
