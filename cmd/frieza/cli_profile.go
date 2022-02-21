package main

import (
	"fmt"
	"log"

	. "github.com/outscale-dev/frieza/internal/common"
	"github.com/teris-io/cli"
)

func cliProfile() cli.Command {
	return cli.NewCommand("profile", "manage cloud profiles").
		WithCommand(cliProfileLs()).
		WithCommand(cliProfileDescribe()).
		WithCommand(cliProfileNew()).
		WithCommand(cliProfileRm()).
		WithCommand(cliProfileTest())
}

func cliProfileLs() cli.Command {
	return cli.NewCommand("list", "list profiles").
		WithShortcut("ls").
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			profileLs(options["config"])
			return 0
		})
}

func cliProfileDescribe() cli.Command {
	return cli.NewCommand("describe", "describe profiles").
		WithShortcut("desc").
		WithArg(cli.NewArg("profile_name", "profile's name to describe")).
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			profileDescribe(options["config"], &args[0])
			return 0
		})
}

func cliProfileNew() cli.Command {
	cmd := cli.NewCommand("new", "create new profile")
	for _, providerCli := range providersCli {
		name, c := providerCli()
		c.WithArg(cli.NewArg("profile", "frieza profile name")).
			WithOption(cliConfigPath()).
			WithOption(cliDebug()).
			WithAction(func(args []string, options map[string]string) int {
				setupDebug(options)
				profile := Profile{
					Name:     args[0],
					Provider: name,
					Config:   options,
				}
				profileNew(options["config"], profile)
				return 0
			})
		cmd.WithCommand(c)
	}
	return cmd
}

func cliProfileRm() cli.Command {
	return cli.NewCommand("remove", "remove profile").
		WithArg(cli.NewArg("profile_name", "profile's name to remove")).
		WithShortcut("rm").
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			profileRm(options["config"], &args[0])
			return 0
		})
}

func cliProfileTest() cli.Command {
	return cli.NewCommand("test", "test profile's authentication").
		WithArg(cli.NewArg("profile_name", "profile's name to test")).
		WithOption(cliConfigPath()).
		WithOption(cliDebug()).
		WithAction(func(args []string, options map[string]string) int {
			setupDebug(options)
			profileTest(options["config"], &args[0])
			return 0
		})
}

func profileLs(customConfigPath string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoad(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}
	for _, profile := range config.Profiles {
		fmt.Println(profile.Name)
	}
}

func profileDescribe(customConfigPath string, profileName *string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoad(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}

	for _, profile := range config.Profiles {
		if *profileName == profile.Name {
			log.Print(profile)
			return
		}
	}
	log.Fatal("Profile not found")
}

func profileNew(customConfigPath string, newProfile Profile) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoad(configPath)
	if err != nil {
		config = ConfigNew()
		if GlobalCliOptions.debug {
			log.Println(err.Error())
		}
	} else {
		for _, profile := range config.Profiles {
			if newProfile.Name == profile.Name {
				log.Fatalf("Profile %s already exist", profile.Name)
			}
		}
	}
	if _, err := ProviderNew(newProfile); err != nil {
		log.Fatalf("Cannot create profile %s: %s", newProfile.Name, err.Error())
	}
	config.Profiles = append(config.Profiles, newProfile)
	config.Write(configPath)
}

func profileRm(customConfigPath string, profileName *string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoad(configPath)
	if err != nil {
		log.Fatal("Cannot load configuration: " + err.Error())
	}

	for idx, profile := range config.Profiles {
		if *profileName == profile.Name {
			config.Profiles = removeProfileIndex(config.Profiles, idx)
			config.Write(configPath)
			return
		}
	}
	log.Fatal("Profile not found")
}

func profileTest(customConfigPath string, profileName *string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoad(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}

	for _, profile := range config.Profiles {
		if *profileName == profile.Name {
			p, err := ProviderNew(profile)
			if err != nil {
				log.Fatalf("Cannot initialize provider %s with profile %s: %s", profile.Provider, profile.Name, err.Error())
			}
			if err := p.AuthTest(); err != nil {
				log.Fatalf("Provider test failed: %s", err.Error())
			}
			return
		}
	}
	log.Fatal("Profile not found")
}

func removeProfileIndex(s []Profile, index int) []Profile {
	return append(s[:index], s[index+1:]...)
}
