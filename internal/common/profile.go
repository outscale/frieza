package common

import (
	"errors"
	"fmt"
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
	out := fmt.Sprintf("profile: %v\n", profile.Name)
	providers, err := profile.GetProviders()
	if err != nil {
		out += fmt.Sprintf("ERROR: %v\n", err)
		return out
	}
	out += "providers:\n"
	for _, provider := range providers {
		out += fmt.Sprintf("  - %v\n", provider)
	}
	out += "configuration:\n"
	for key, value := range profile.Config {
		out += fmt.Sprintf("  - %v: %v\n", key, value)
	}
	return out
}
