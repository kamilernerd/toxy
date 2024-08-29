package toxy

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	Random     string = "random"
	Sequential        = "sequential"
)

const (
	ServiceUp      string = "up"
	ServiceDown    string = "down"
	ServiceUnknown string = "unknown"
)

type ServerConfig struct {
	Port     int    `toml:"port"`
	Hostname string `toml:"hostname"`
	Name     string `toml:"name"`
}

type Config struct {
	Port            int    `toml:"port"`
	Hostname        string `toml:"hostname"`
	CertPath        string `toml:"cert_file"`
	KeyPath         string `toml:"key_file"`
	LoadBalancer    string `toml:"load_balancer"`
	ResolveInterval int    `toml:"resolve_interval"`
	Server          map[string][]ServerConfig
}

func LoadConfig() Config {
	content, err := os.ReadFile("./config.toml")
	if err != nil {
		log.Fatalf("Failed to load server config \n %v", err)
	}

	defaultConfStruct := Config{}

	_, err = toml.Decode(string(content), &defaultConfStruct)
	if err != nil {
		log.Fatalf("Failed to parse server config \n %v", err)
	}

	return defaultConfStruct
}
