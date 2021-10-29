package main

import (
	"fmt"

	. "github.com/outscale-dev/frieza/internal/common"
	"github.com/outscale-dev/frieza/internal/providers/outscale_oapi"
	"github.com/outscale-dev/frieza/internal/providers/provider_example"
	"github.com/teris-io/cli"
)

func ProviderNew(profile Profile) (Provider, error) {
	switch profile.Provider {
	case outscale_oapi.Name:
		return outscale_oapi.New(profile.Config, Debug)
	case provider_example.Name:
		return provider_example.New(profile.Config, Debug)
	}
	return nil, fmt.Errorf("Provider %s not found", profile.Provider)
}

var providersCli = []func() (string, cli.Command){
	outscale_oapi.Cli,
	//provider_example.Cli,
}

var providersTypes = map[string][]ObjectType{
	outscale_oapi.Name:    outscale_oapi.Types(),
	provider_example.Name: provider_example.Types(),
}
