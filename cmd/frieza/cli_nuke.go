package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	. "github.com/outscale-dev/frieza/internal/common"
	"github.com/teris-io/cli"
)

func cliNuke() cli.Command {
	return cli.NewCommand("nuke", "delete ALL resources of specified profiles").
		WithOption(cli.NewOption("plan", "Only show what resource would be deleted").WithType(cli.TypeBool)).
		WithOption(cli.NewOption("auto-approve", "Approve resource deletion without confirmation").WithType(cli.TypeBool)).
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithArg(cli.NewArg("profile", "one or more profile").AsOptional()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			nuke(options["config"], args, options["plan"] == "true", options["auto-approve"] == "true")
			return 0
		})
}

func nuke(customConfigPath string, profiles []string, plan bool, autoApprove bool) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}

	if len(profiles) == 0 {
		log.Fatalln("No profile provided, use --help for more details.")
	}

	uniqueProfiles := make(map[string]bool)
	for _, profile := range profiles {
		uniqueProfiles[profile] = true
	}

	config, err := ConfigLoad(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}

	var providers []Provider
	var objectsToDelete []Objects
	objectsCount := 0

	for profileName := range uniqueProfiles {
		profile, err := config.GetProfile(profileName)
		if err != nil {
			log.Fatalf("Error while getting profile %s: %s", profileName, err.Error())
		}
		provider, err := ProviderNew(*profile)
		if err != nil {
			log.Fatalf("Error intializing profile %s: %s", profileName, err.Error())
		}
		toDelete := ReadObjects(&provider)
		objectsCount += ObjectsCount(&toDelete)
		fmt.Printf("Profile %s (%s):\n", profile.Name, provider.Name())
		if objectsCount > 0 {
			fmt.Print(ObjectsPrint(&provider, &toDelete))
		} else {
			fmt.Println("* no object *")
		}
		providers = append(providers, provider)
		objectsToDelete = append(objectsToDelete, toDelete)
	}

	if objectsCount == 0 {
		fmt.Println("\nNothing to delete, exiting")
		return
	}

	if plan {
		return
	}

	message := fmt.Sprintf("Do you really want to delete ALL resources?\n" +
		"  Frieza will delete all resources shown above.")
	if !confirmAction(&message, autoApprove) {
		log.Fatal("Nuke canceled")
	}
	loopDelete(providers, objectsToDelete)
}

func loopDelete(providers []Provider, objects []Objects) {
	for {
		var objectsCount []int
		var totalObjectCount int
		for i := range objects {
			count := ObjectsCount(&objects[i])
			totalObjectCount += count
			objectsCount = append(objectsCount, count)
		}
		if totalObjectCount == 0 {
			return
		}
		for i, provider := range providers {
			if objectsCount[i] == 0 {
				continue
			}
			DeleteObjects(&provider, objects[i])
		}
		for i, provider := range providers {
			diff := NewDiff()
			remaining := ReadNonEmptyObjects(&provider, objects[i])
			diff.Build(&remaining, &objects[i])
			objects[i] = diff.Retained
		}
	}
}

func confirmAction(message *string, autoApprove bool) bool {
	if autoApprove {
		return true
	}
	fmt.Printf("\n%s\n", *message)
	fmt.Printf("  There is no undo. Only 'yes' will be accepted to confirm.\n\n")
	fmt.Printf("  Enter a value: ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.Replace(response, "\n", "", -1)
	response = strings.Replace(response, "\r", "", -1)
	if response != "yes" {
		return false
	}
	return true
}
