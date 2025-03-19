package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/BrandonBentley/ezbind"
	"github.com/spf13/viper"
)

// Emulating a file set to be read
var jsonConfigFile = strings.NewReader(`{"server":{"http":{"port":80}}}`)

func init() {
	// Assume these are set before running
	os.Setenv("AUTH_TOKEN", "ABCDEFG")
}

type Config struct {
	Server Server `mapstructure:"server"`
	Auth   Auth   `mapstructure:"auth"`
}

type Server struct {
	Http Http `mapstructure:"http"`
}

type Http struct {
	Port int `mapstructure:"port"`
}

type Auth struct {
	Token string `mapstructure:"token"`
}

func main() {
	cfg := viper.New()
	cfg.SetConfigType("json")

	var config Config

	ezbind.BindStruct(cfg, config)

	cfg.ReadConfig(jsonConfigFile)

	cfg.Unmarshal(&config)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	enc.Encode(config)
}
