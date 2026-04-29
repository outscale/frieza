package common

import (
	"errors"
	"fmt"
	"strings"
)

type Profile struct {
	Name      string         `json:"name"`
	Provider  string         `json:"provider,omitempty"`
	Providers []string       `json:"providers,omitempty"`
	Config    ProviderConfig `json:"config"`
}

func (p *Profile) GetProviders() ([]string, error) {
	if p.Provider != "" && len(p.Providers) > 0 {
		return nil, errors.New("profile has both 'provider' and 'providers' fields specified; use 'providers' instead")
	}

	if len(p.Providers) > 0 {
		return p.Providers, nil
	}

	if p.Provider != "" {
		return []string{p.Provider}, nil
	}

	return []string{}, nil
}

func (profile Profile) String() string {
	var outBuilder strings.Builder

	fmt.Fprintf(&outBuilder, "profile: %v\n", profile.Name)
	providers, err := profile.GetProviders()
	if err != nil {
		fmt.Fprintf(&outBuilder, "ERROR: %v\n", err)
		return outBuilder.String()
	}

	outBuilder.WriteString("providers:\n")
	for _, provider := range providers {
		fmt.Fprintf(&outBuilder, "  - %v\n", provider)
	}

	outBuilder.WriteString("configuration:\n")
	for key, value := range profile.Config {
		fmt.Fprintf(&outBuilder, "  - %v: %v\n", key, value)
	}
	return outBuilder.String()
}
