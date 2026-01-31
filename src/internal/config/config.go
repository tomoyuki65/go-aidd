package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Issue struct {
		Provider string `koanf:"provider"`
		Label    string `koanf:"label"`
	} `koanf:"issue"`
	GitHub struct {
		Repository           string `koanf:"repository"`
		CloneType            string `koanf:"clone_type"`
		CloneBranch          string `koanf:"clone_branch"`
		PushBranchOnComplete bool   `koanf:"push_branch_on_complete"`
		CreatePrOnComplete   bool   `koanf:"create_pr_on_complete"`
		PrDraft              bool   `koanf:"pr_draft"`
	} `koanf:"github"`
	Task struct {
		ListPageSize int `koanf:"list_page_size"`
		SkipRunTask bool `koanf:"skip_run_task"`
	} `koanf:"task"`
	AI struct {
		Type  string `koanf:"type"`
		Model string `koanf:"model"`
	} `koanf:"ai"`
}

func getConfigPath() string {
	searchPaths := []string{
		"config.yml",
		filepath.Join("src", "config.yml"),
	}

	var configPath string
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	if configPath == "" {
		log.Fatal("no config.yml found in search paths")
	}

	return configPath
}

func LoadConfig() *Config {
	// Load configuration from config.yml
	configPath := getConfigPath()

	k := koanf.New(".")
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return &cfg
}
