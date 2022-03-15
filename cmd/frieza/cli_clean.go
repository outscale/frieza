package main

import (
	"fmt"
	"log"

	. "github.com/outscale-dev/frieza/internal/common"
	"github.com/teris-io/cli"
)

func cliClean() cli.Command {
	return cli.NewCommand("clean", "delete created resources since a specific snapshot").
		WithOption(cli.NewOption("plan", "Only show what resource would be deleted").WithType(cli.TypeBool)).
		WithOption(cli.NewOption("auto-approve", "Approve resource deletion without confirmation").WithType(cli.TypeBool)).
		WithArg(cli.NewArg("snapshot_name", "snapshot")).
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			clean(options["config"], &args[0], options["plan"] == "true", options["auto-approve"] == "true")
			return 0
		})
}

func clean(customConfigPath string, snapshotName *string, plan bool, autoApprove bool) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}
	snapshot, err := SnapshotLoad(*snapshotName, config.SnapshotFolderPath)
	if err != nil {
		log.Fatalf("Error load snapshot %s: %s", *snapshotName, err.Error())
	}

	var providers []Provider
	var objectsToDelete []Objects
	objectsCount := 0

	for _, data := range snapshot.Data {
		profile, err := config.GetProfile(data.Profile)
		if err != nil {
			log.Fatalf("Error while getting profile %s: %s", data.Profile, err.Error())
		}
		provider, err := ProviderNew(*profile)
		if err != nil {
			log.Fatalf("Error intializing profile %s: %s", data.Profile, err.Error())
		}
		objects := ReadObjects(&provider)
		diff := NewDiff()
		diff.Build(&data.Objects, &objects)
		count := ObjectsCount(&diff.Created)
		objectsCount += count
		if count > 0 {
			fmt.Printf("Newly created object to delete in profile %s (%s):\n", profile.Name, provider.Name())
			fmt.Printf(ObjectsPrint(&provider, &diff.Created))
		} else {
			fmt.Printf("No new object to delete in profile %s (%s)\n", profile.Name, provider.Name())
		}
		providers = append(providers, provider)
		objectsToDelete = append(objectsToDelete, *&diff.Created)
	}

	if objectsCount == 0 {
		fmt.Println("Nothing to delete, exiting")
		return
	}

	if plan {
		return
	}

	message := fmt.Sprintf("Do you really want to delete newly created resources?\n" +
		"  Frieza will delete all resources shown above.")
	if !confirmAction(&message, autoApprove) {
		log.Fatal("Clean canceled")
	}
	loopDelete(providers, objectsToDelete)
}
