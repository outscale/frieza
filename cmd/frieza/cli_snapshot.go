package main

import (
	"log"
	"os"
	"strings"
	"time"

	. "github.com/outscale/frieza/internal/common"
	"github.com/teris-io/cli"
)

func cliSnapshot() cli.Command {
	return cli.NewCommand("snapshot", "manage resource snapshots").
		WithShortcut("snap").
		WithCommand(cliSnapshotNew()).
		WithCommand(cliSnapshotLs()).
		WithCommand(cliSnapshotDescribe()).
		WithCommand(cliSnapshotRm()).
		WithCommand(cliSnapshotUpdate())
}

func cliSnapshotNew() cli.Command {
	return cli.NewCommand("new", "create new snapshot containing all resource ids").
		WithArg(cli.NewArg("snapshot_name", "snapshot name")).
		WithOption(cli.NewOption("only-resource-types", "Remove only theses resource types (separated by ','). You can see all resource types in the description of the provider.").WithType(cli.TypeString)).
		WithOption(cli.NewOption("exclude-resource-types", "Remove all except theses resource types (separated by ','). You can see all resource types in the description of the provider.").WithType(cli.TypeString)).
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithArg(cli.NewArg("profile", "one or more profile to snapshot").AsOptional()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			snapshotNew(options["config"], args, options)
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

func cliSnapshotUpdate() cli.Command {
	return cli.NewCommand("update", "update existing snapshot").
		WithShortcut("up").
		WithArg(cli.NewArg("snapshot_name", "snapshot name")).
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithOption(cli.NewOption("incremental", "update snapshot incrementally").WithType(cli.TypeBool).WithChar('i')).
		WithAction(func(args []string, options map[string]string) int {
			incrementalUpdate := options["incremental"] == "true"
			setupDebug(options)
			snapshotUpdate(options["config"], &args[0], incrementalUpdate)
			return 0
		})
}

func snapshotNew(customConfigPath string, args []string, options map[string]string) {
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

	var providers []Provider
	var profiles []string
	for _, profileName := range profileNames {
		found := false
		for _, profile := range config.Profiles {
			if profileName == profile.Name {
				profileProviders, err := ProviderNew(profile)
				if err != nil {
					log.Fatalf("Cannot initialize profile %s: %s", profile.Name, err.Error())
				}
				providers = append(providers, profileProviders...)
				for range profileProviders {
					profiles = append(profiles, profile.Name)
				}
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

	date := time.Now().UTC().String()
	snapshot := Snapshot{
		Version: SnapshotVersion(),
		Name:    snapshotName,
		Date:    date,
		Config:  config,
		Filters: resourcesTypeFilterPtr,
	}
	for i, provider := range providers {
		objs, err := ReadObjects(&provider, resourcesTypeFilterPtr)
		if err != nil {
			log.Fatalf("Error reading objects: %v\n", err)
		}
		snapshot.Data = append(snapshot.Data, SnapshotData{
			Profile:  profiles[i],
			Provider: provider.Name(),
			Objects:  objs,
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
	files, err := os.ReadDir(config.SnapshotFolderPath)
	if err != nil {
		log.Fatalf("Error while listing snapshots: %s", err.Error())
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		snapshotName := strings.TrimSuffix(file.Name(), ".json")
		if snapshot, err := SnapshotLoad(snapshotName, config); err == nil {
			log.Println(snapshot.Name)
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
	snapshot, err := SnapshotLoad(*snapshotName, config)
	if err != nil {
		log.Fatalf("Cannot load snapshot %s: %s", *snapshotName, err.Error())
	}
	log.Print(snapshot)
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

func snapshotUpdate(customConfigPath string, snapshotName *string, incrementalUpdate bool) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}
	var snapshot *Snapshot
	if snapshot, err = SnapshotLoad(*snapshotName, config); err != nil {
		log.Fatalf("Snapshot %s does not exist", *snapshotName)
	}

	for _, data := range snapshot.Data {
		profile, err := config.GetProfile(data.Profile)
		if err != nil {
			log.Fatalf("Error while getting profile %s: %s", data.Profile, err.Error())
		}
		providers, err := ProviderNew(*profile)
		if err != nil {
			log.Fatalf("Error intializing profile %s: %s", data.Profile, err.Error())
		}

		for _, provider := range providers {
			if err := provider.AuthTest(); err != nil {
				log.Fatalf("Provider %s test failed for profile %s: %s", provider.Name(), profile.Name, err.Error())
			}

			objects, err := ReadObjects(&provider, snapshot.Filters)
			if err != nil {
				log.Fatalf("Error reading objects: %v\n", err)
			}
			diff := NewDiff()
			diff.Build(&data.Objects, &objects)
			for key, value := range diff.Created {
				var objectToAdd []string
				if incrementalUpdate {
					incrementObject, err := incrementalChoice(key, value)
					if err != nil {
						log.Fatalf("Snapshot failed: %s", err.Error())
					}

					if incrementObject == nil {
						log.Fatalf("Snapshot update cancels")
					}
					objectToAdd = append(objectToAdd, (*incrementObject)...)
				} else {
					objectToAdd = value
				}

				if snapshotValue, ok := data.Objects[key]; ok {
					data.Objects[key] = append(snapshotValue, objectToAdd...)
				} else {
					data.Objects[key] = objectToAdd
				}

			}
		}
	}

	date := time.Now().UTC().String()
	snapshot.Version = SnapshotVersion()
	snapshot.Date = date

	if err = snapshot.Write(); err != nil {
		log.Fatalf("Snapshot failed: %s", err.Error())
	}
}
