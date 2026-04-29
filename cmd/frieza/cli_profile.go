package main

import (
	"context"
	"log"
	"slices"

	. "github.com/outscale/frieza/internal/common"
	"github.com/teris-io/cli"
)

func cliProfile() cli.Command {
	return cli.NewCommand("profile", "manage cloud profiles").
		WithCommand(cliProfileLs()).
		WithCommand(cliProfileDescribe()).
		WithCommand(cliProfileNew()).
		WithCommand(cliProfileRm()).
		WithCommand(cliProfileTest()).
		WithCommand(cliProfileAddProvider()).
		WithCommand(cliProfileRemoveProvider())
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
					Name:      args[0],
					Providers: []string{name},
					Config:    options,
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
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}
	for _, profile := range config.Profiles {
		log.Println(profile.Name)
	}
}

func profileDescribe(customConfigPath string, profileName *string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
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
	config, err := ConfigLoadWithDefault(configPath)
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
	_ = config.Write(configPath)
}

func profileRm(customConfigPath string, profileName *string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatal("Cannot load configuration: " + err.Error())
	}

	for idx, profile := range config.Profiles {
		if *profileName == profile.Name {
			config.Profiles = removeProfileIndex(config.Profiles, idx)
			_ = config.Write(configPath)
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
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}

	ctx := context.Background()

	for _, profile := range config.Profiles {
		if *profileName == profile.Name {
			providers, err := ProviderNew(profile)
			if err != nil {
				log.Fatalf(
					"Cannot initialize providers for profile %s: %s",
					profile.Name,
					err.Error(),
				)
			}
			for _, p := range providers {
				log.Printf("Testing provider %s...", p.Name())
				if err := p.AuthTest(ctx); err != nil {
					log.Fatalf("Provider %s test failed: %s", p.Name(), err.Error())
				}
				log.Printf("Provider %s test passed", p.Name())
			}
			return
		}
	}
	log.Fatal("Profile not found")
}

func removeProfileIndex(s []Profile, index int) []Profile {
	return append(s[:index], s[index+1:]...)
}

func cliProfileAddProvider() cli.Command {
	cmd := cli.NewCommand("add-provider", "add a provider to an existing profile")

	for _, providerCli := range providersCli {
		providerName, _ := providerCli()
		subCmd := cli.NewCommand(providerName, "add "+providerName+" provider").
			WithArg(cli.NewArg("profile_name", "profile's name")).
			WithOption(cliConfigPath()).
			WithOption(cliDebug()).
			WithAction(func(args []string, options map[string]string) int {
				setupDebug(options)
				profileName := args[0]
				profileAddProvider(options["config"], profileName, providerName)
				return 0
			})
		cmd.WithCommand(subCmd)
	}
	return cmd
}

func cliProfileRemoveProvider() cli.Command {
	cmd := cli.NewCommand("remove-provider", "remove a provider from a profile")

	for _, providerCli := range providersCli {
		providerName, _ := providerCli()
		subCmd := cli.NewCommand(providerName, "remove "+providerName+" provider").
			WithArg(cli.NewArg("profile_name", "profile's name")).
			WithOption(cliConfigPath()).
			WithOption(cliDebug()).
			WithAction(func(args []string, options map[string]string) int {
				setupDebug(options)
				profileName := args[0]
				profileRemoveProvider(options["config"], profileName, providerName)
				return 0
			})
		cmd.WithCommand(subCmd)
	}
	return cmd
}

func profileAddProvider(customConfigPath string, profileName string, providerName string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}

	for i := range config.Profiles {
		profile := &config.Profiles[i]
		if profileName == profile.Name {
			currentProviders, err := profile.GetProviders()
			if err != nil {
				log.Fatalf("Error reading profile providers: %s", err.Error())
			}

			if slices.Contains(currentProviders, providerName) {
				log.Fatalf("Provider %s already exists in profile %s", providerName, profileName)
			}
			profile.Provider = ""
			profile.Providers = append(currentProviders, providerName)
			if err := config.Write(configPath); err != nil {
				log.Fatalf("Failed to write config: %s", err.Error())
			}
			log.Printf("Added provider %s to profile %s", providerName, profileName)
			return
		}
	}
	log.Fatal("Profile not found")
}

func profileRemoveProvider(customConfigPath string, profileName string, providerName string) {
	var configPath *string
	if len(customConfigPath) > 0 {
		configPath = &customConfigPath
	}
	config, err := ConfigLoadWithDefault(configPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err.Error())
	}

	for i := range config.Profiles {
		profile := &config.Profiles[i]
		if profileName == profile.Name {
			currentProviders, err := profile.GetProviders()
			if err != nil {
				log.Fatalf("Error reading profile providers: %s", err.Error())
			}

			newProviders := slices.DeleteFunc(currentProviders, func(p string) bool { return p == providerName })
			if len(newProviders) == len(currentProviders) {
				log.Fatalf("Provider %s not found in profile %s", providerName, profileName)
			}
			if len(newProviders) == 0 {
				log.Fatalf("Cannot remove the last provider from profile %s. Delete the profile instead.", profileName)
			}

			profile.Provider = ""
			profile.Providers = newProviders
			if err := config.Write(configPath); err != nil {
				log.Fatalf("Failed to write config: %s", err.Error())
			}
			log.Printf("Removed provider %s from profile %s", providerName, profileName)
			return
		}
	}
	log.Fatal("Profile not found")
}
