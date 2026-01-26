package config

import (
	_ "embed"
	"log"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

//go:embed config.yml
var configYAML []byte

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

func LoadConfig() *Config {
	// Load configuration from config.yml
	k := koanf.New(".")
	if err := k.Load(rawbytes.Provider(configYAML), yaml.Parser()); err != nil {
		log.Fatalf("failed to load embedded config: %v", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return &cfg
}
