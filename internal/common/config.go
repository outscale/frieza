package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var version string
var commit string

func FullVersion() string {
	if len(version) == 0 {
		version = "0.0.0-beta-unknown-version"
	}
	if len(commit) == 0 {
		commit = "unknown git commit"
	}
	return fmt.Sprintf("%s-%s", version, commit)
}

func ConfigVersion() int {
	return 0
}

type Config struct {
	Version            int       `json:"version"`
	Profiles           []Profile `json:"profiles"`
	SnapshotFolderPath string    `json:"snapshot_folder_path,omitempty"`
}

func ConfigNew() *Config {
	return &Config{
		Version:            ConfigVersion(),
		Profiles:           []Profile{},
		SnapshotFolderPath: "",
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

func DefaultSnapshotFolderPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".frieza", "snapshots"), nil
}

func ConfigLoadWithDefault(customConfigPath *string) (*Config, error) {
	config, err := ConfigLoad(customConfigPath)
	if err != nil {
		return nil, err
	}
	if len(config.SnapshotFolderPath) == 0 {
		config.SnapshotFolderPath, err = DefaultSnapshotFolderPath()
		if err != nil {
			return nil, err
		}
	}
	return config, nil
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
