package main

import (
	"fmt"

	. "github.com/outscale/frieza/internal/common"
	"github.com/outscale/frieza/internal/providers/fs"
	oapi "github.com/outscale/frieza/internal/providers/outscale_oapi"
	oks "github.com/outscale/frieza/internal/providers/outscale_oks"
	"github.com/outscale/frieza/internal/providers/s3"
	"github.com/teris-io/cli"
)

func ProviderNew(profile Profile) ([]Provider, error) {
	var providers []Provider
	providerNames, err := profile.GetProviders()
	if err != nil {
		return nil, err
	}
	for _, providerName := range providerNames {
		var provider Provider
		var err error
		switch providerName {
		case oapi.Name:
			provider, err = oapi.New(profile.Config, GlobalCliOptions.debug)
		case s3.Name:
			provider, err = s3.New(profile.Config, GlobalCliOptions.debug)
		case fs.Name:
			provider, err = fs.New(profile.Config, GlobalCliOptions.debug)
		case oks.Name:
			provider, err = oks.New(profile.Config, GlobalCliOptions.debug)
		default:
			return nil, fmt.Errorf("provider %s not found", providerName)
		}
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	if len(providers) == 0 {
		return nil, fmt.Errorf("no providers configured for profile %s", profile.Name)
	}
	return providers, nil
}

var providersCli = []func() (string, cli.Command){
	oapi.Cli,
	s3.Cli,
	fs.Cli,
	oks.Cli,
}

var providersTypes = map[string][]ObjectType{
	oapi.Name: oapi.Types(),
	oks.Name:  oks.Types(),
	s3.Name:   s3.Types(),
	fs.Name:   fs.Types(),
}
