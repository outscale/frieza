package common

import "fmt"

type Profile struct {
	Name     string         `json:"name"`
	Provider string         `json:"provider"`
	Config   ProviderConfig `json:"config"`
}

func (profile Profile) String() string {
	out := fmt.Sprintf("profile: %v\n", profile.Name)
	out += fmt.Sprintf("provider: %v\n", profile.Provider)
	out += "configuration:\n"
	for key, value := range profile.Config {
		out += fmt.Sprintf("  - %v: %v\n", key, value)
	}
	return out
}
