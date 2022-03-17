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
		WithOption(cliJson()).
		WithOption(cli.NewOption("auto-approve", "Approve resource deletion without confirmation").WithType(cli.TypeBool)).
		WithArg(cli.NewArg("snapshot_name", "snapshot")).
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			clean(options["config"], &args[0], options["plan"] == "true", options["auto-approve"] == "true", options["json"] == "true")
			return 0
		})
}

func clean(customConfigPath string, snapshotName *string, plan bool, autoApprove bool, jsonOutput bool) {
	var configPath *string
	if jsonOutput && !autoApprove {
		cliFatalf(true, "Cannot use --json option without --auto-approve")
	}
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		cliFatalf(jsonOutput, "Cannot load configuration: %s", err.Error())
	}
	snapshot, err := SnapshotLoad(*snapshotName, config)
	if err != nil {
		cliFatalf(jsonOutput, "Error load snapshot %s: %s", *snapshotName, err.Error())
	}

	destroyer := NewDestroyer()
	objectsCount := 0
	for _, data := range snapshot.Data {
		profile, err := config.GetProfile(data.Profile)
		if err != nil {
			cliFatalf(jsonOutput, "Error while getting profile %s: %s", data.Profile, err.Error())
		}
		provider, err := ProviderNew(*profile)
		if err != nil {
			cliFatalf(jsonOutput, "Error intializing profile %s: %s", data.Profile, err.Error())
		}
		objects := ReadObjects(&provider)
		diff := NewDiff()
		diff.Build(&data.Objects, &objects)
		objectsCount += ObjectsCount(&diff.Created)
		destroyer.add(profile, &provider, &diff.Created)
	}

	destroyer.print(jsonOutput)
	if plan || objectsCount == 0 {
		return
	}
	if jsonOutput {
		disableLogs()
	}
	message := fmt.Sprintf("Do you really want to delete newly created resources?\n" +
		"  Frieza will delete all resources shown above.")
	if !confirmAction(&message, autoApprove) {
		log.Fatal("Clean canceled")
	}
	destroyer.run()
}
