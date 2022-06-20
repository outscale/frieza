package main

import (
	"fmt"

	"github.com/teris-io/cli"
)

func cliProvider() cli.Command {
	return cli.NewCommand("provider", "show supported providers and their features").
		WithCommand(cliProviderLs()).
		WithCommand(cliProviderDescribe())
}

func cliProviderLs() cli.Command {
	return cli.NewCommand("list", "list providers").
		WithShortcut("ls").
		WithAction(func(args []string, options map[string]string) int {
			for providerName := range providersTypes {
				fmt.Println(providerName)
			}
			return 0
		})
}

func cliProviderDescribe() cli.Command {
	return cli.NewCommand("describe", "describe provider features").
		WithShortcut("desc").
		WithArg(cli.NewArg("provider_name", "provider to describe")).
		WithAction(func(args []string, options map[string]string) int {
			providerName := args[0]
			for _, providerType := range providersTypes[providerName] {
				fmt.Println(providerType)
			}
			return 0
		})
}
