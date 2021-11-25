package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func ConfigVersion() int {
	return 0
}

type Config struct {
	Version  int       `json:"version"`
	Profiles []Profile `json:"profiles"`
}

func ConfigNew() *Config {
	return &Config{
		Version:  ConfigVersion(),
		Profiles: []Profile{},
	}
}

func DefaultConfigFolderPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".frieza"), nil
}

func DefaultConfigPath() (string, error) {
	folderPath, err := DefaultConfigFolderPath()
	if err != nil {
		return "", err
	}
	return path.Join(folderPath, "config.json"), nil
}

func ConfigLoad(customConfigPath *string) (*Config, error) {
	var configPath string
	var err error
	if customConfigPath == nil {
		configPath, err = DefaultConfigPath()
		if err != nil {
			return nil, err
		}
	} else {
		configPath = *customConfigPath
	}
	config_json, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(config_json, &config); err != nil {
		return nil, err
	}
	if config.Version > ConfigVersion() {
		return nil, errors.New("configuration version not supported, please upgrade frieza")
	}
	return &config, nil
}

func (config *Config) Write(customConfigPath *string) error {
	var configPath string
	if customConfigPath == nil {
		configFolderPath, err := DefaultConfigFolderPath()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(configFolderPath, os.ModePerm); err != nil {
			return err
		}
		configPath, err = DefaultConfigPath()
		if err != nil {
			return err
		}
	} else {
		configPath = *customConfigPath
	}
	json_bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(configPath, json_bytes, 0700); err != nil {
		return err
	}
	return nil
}

func (config *Config) GetProfile(profileName string) (*Profile, error) {
	for _, profile := range config.Profiles {
		if profileName == profile.Name {
			return &profile, nil
		}
	}
	return nil, fmt.Errorf("Profile %s not found", profileName)
}
