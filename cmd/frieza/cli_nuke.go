package main

import (
	"context"
	"log"
	"strings"
	"time"

	. "github.com/outscale/frieza/internal/common"
	"github.com/teris-io/cli"
)

func cliNuke() cli.Command {
	return cli.NewCommand("nuke", "delete ALL resources of specified profiles").
		WithOption(cli.NewOption("plan", "Only show what resource would be deleted").WithType(cli.TypeBool)).
		WithOption(cli.NewOption("timeout", "Exit with error after a specific duration (ex: 30s, 5m, 1.5h)").WithType(cli.TypeString)).
		WithOption(cli.NewOption("only-resource-types", "Remove only theses resource types (separated by ','). You can see all resource types in the description of the provider.").WithType(cli.TypeString)).
		WithOption(cli.NewOption("exclude-resource-types", "Remove all except theses resource types (separated by ','). You can see all resource types in the description of the provider.").WithType(cli.TypeString)).
		WithOption(cliJson()).
		WithOption(cli.NewOption("auto-approve", "Approve resource deletion without confirmation").WithType(cli.TypeBool)).
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithArg(cli.NewArg("profile", "one or more profile").AsOptional()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			plan := options["plan"] == "true"
			autoApprove := options["auto-approve"] == "true"
			jsonOutput := options["json"] == "true"
			timeout := "-1"
			if len(options["timeout"]) > 0 {
				timeout = options["timeout"]
			}

			var resourcesTypeFilterPtr *ResourceFilterEnvelope
			if len(options["only-resource-types"]) > 0 && len(options["exclude-resource-types"]) > 0 {
				cliFatalf(true, "Cannot use --only-resource-types option with --exclude-resource-types")
			}
			if len(options["only-resource-types"]) > 0 {
				resourcesTypeFilterPtr = NewResourceFilterOnly(strings.Split(options["only-resource-types"], ","))
			}
			if len(options["exclude-resource-types"]) > 0 {
				resourcesTypeFilterPtr = NewResourceFilterExclude(strings.Split(options["exclude-resource-types"], ","))
			}

			nuke(options["config"], args, plan, autoApprove, jsonOutput, timeout, resourcesTypeFilterPtr)
			return 0
		})
}

func nuke(customConfigPath string, profiles []string, plan bool, autoApprove bool, jsonOutput bool, timeout string, resourceFilter *ResourceFilterEnvelope) {
	if jsonOutput && !autoApprove {
		cliFatalf(true, "Cannot use --json option without --auto-approve")
	}
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}

	if len(profiles) == 0 {
		cliFatalf(jsonOutput, "No profile provided, use --help for more details.")
	}

	uniqueProfiles := make(map[string]bool)
	for _, profile := range profiles {
		uniqueProfiles[profile] = true
	}

	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		cliFatalf(jsonOutput, "Cannot load configuration: %s", err.Error())
	}

	ctx := context.Background()

	destroyer := NewDestroyer()
	for profileName := range uniqueProfiles {
		profile, err := config.GetProfile(profileName)
		if err != nil {
			cliFatalf(jsonOutput, "Error while getting profile %s: %s", profileName, err.Error())
		}
		providers, err := ProviderNew(*profile)
		if err != nil {
			cliFatalf(jsonOutput, "Error intializing profile %s: %s", profileName, err.Error())
		}
		for _, provider := range providers {
			objectsToDelete, err := ReadObjects(ctx, &provider, resourceFilter)
			if err != nil {
				log.Fatalf("Error reading objects: %v", err)
			}
			destroyer.add(profile, &provider, &objectsToDelete)
		}
	}

	destroyer.print(jsonOutput)
	if plan {
		return
	}
	if jsonOutput {
		disableLogs()
	}
	message := "Do you really want to delete ALL resources?\n" +
		"  Frieza will delete all resources shown above."
	if !confirmAction(&message, autoApprove) {
		log.Fatal("Nuke canceled")
	}

	tout, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatal("Could not parse timeout: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, tout)
	defer cancel()

	destroyer.run(ctx)
}
