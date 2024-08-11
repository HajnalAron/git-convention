package config

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hajnalaron/git-convention-cli/types"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultBranchPrefix string         `json:"default_branch_prefix"`
	DefaultCommitPrefix string         `json:"default_commit_prefix"`
	EmojisEnabled       bool           `json:"emojis_enabled"`
	BranchTypes         []types.Branch `json:"branch_types"`
	CommitTypes         []types.Commit `json:"commit_types"`
}

//go:embed default.json
var defaultConfig embed.FS

func loadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetConfig(configPathArg string) (*Config, error) {
	fmt.Println(configPathArg)
	configPath := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	var config *Config
	var err error

	if *configPath != "" {
		config, err = loadConfig(*configPath)
		if err != nil {
			log.Fatalf("Error reading config file from provided path: %v", err)
		}

		return config, nil
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Error getting home directory: %v", err)
		}

		defaultConfigPath := filepath.Join(homeDir, ".config", "git-convention-cli", "config.json")

		if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
			fmt.Println("Creating config file in home directory...")
			err = createDefaultConfig(defaultConfigPath)
			if err != nil {
				log.Fatalf("Error creating default config file: %v", err)
			}
			fmt.Println("Config file created")
			config, err = loadConfig(defaultConfigPath)
			if err != nil {
				log.Fatalf("Error reading config file from home directory: %v", err)
			}

			return config, nil
		} else {
			config, err = loadConfig(defaultConfigPath)
			if err != nil {
				log.Fatalf("Error reading config file from home directory: %v", err)
			}
			return config, nil
		}
	}

}

func createDefaultConfig(configPath string) error {
	data, err := defaultConfig.ReadFile("default.json")
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ShowConfig(config *Config) {
	fmt.Printf("Default Branch Prefix: %s\n", config.DefaultBranchPrefix)
	fmt.Printf("Default Commit Prefix: %s\n", config.DefaultCommitPrefix)
	fmt.Printf("Emojis Enabled: %v\n", config.EmojisEnabled)

	fmt.Println("Branches:\n---------")
	for _, branchType := range config.BranchTypes {
		fmt.Printf("	%s: %s \n",
			branchType.Type, branchType.Description)
	}

	fmt.Println("Commits:\n--------")
	for _, commitType := range config.CommitTypes {
		fmt.Printf("	%s: %s %s \n",
			commitType.Type, commitType.Emoji, commitType.Description)
	}
}
