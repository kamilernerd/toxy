package toxy_test

import (
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/kamilernerd/toxy"
)

func TestLoadConfigFile(t *testing.T) {
	_, err := os.ReadFile("./config.toml")
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseConfig(t *testing.T) {
	content, err := os.ReadFile("./config.toml")
	if err != nil {
		t.Fatal(err)
	}

	defaultConfStruct := toxy.Config{}

	_, err = toml.Decode(string(content), &defaultConfStruct)
	if err != nil {
		t.Fatal(err)
	}
}
