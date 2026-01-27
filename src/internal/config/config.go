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
		Repository string `koanf:"repository"`
	} `koanf:"github"`
	AI struct {
		Model string `koanf:"model"`
	} `koanf:"ai"`
}

func getConfigPath() string {
	// If the ENV environment variable is set to "local"
	if os.Getenv("ENV") == "local" {
		return "config.yml"
	}

	// In case of binary execution
	exePath, _ := os.Executable()
	binDir := filepath.Dir(exePath)
	return filepath.Join(binDir, "..", "config.yml")
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
