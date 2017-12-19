package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var config *viper.Viper

const (
	appName = "kongctl"
)

func GetConfig() *viper.Viper {
	return config
}

func Init() error {
	config = viper.New()

	setup()

	err := config.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		return errors.Wrap(err, "Fatal error loading config file")
	}

	setDefaults()

	return nil
}

func setup() {
	// Load in YAML files
	config.SetConfigType("yaml")

	// Allow us to set config using Environment Variables
	config.SetEnvPrefix("")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Want to make sure we're loading in the config file at
	// the right paths.
	config.SetConfigName(fmt.Sprintf("%s", appName))        // name of config file (without extension)
	config.AddConfigPath(fmt.Sprintf("/etc/%s/", appName))  // path to look for the config file in
	config.AddConfigPath(fmt.Sprintf("$HOME/.%s", appName)) // call multiple times to add many search paths
	config.AddConfigPath(".")                               // optionally look for config in the working directory
}

func setDefaults() {
	// By default we'll assume that we're running the Kong
	// service on the same host machine.
	config.SetDefault("host", "http://127.0.0.1:8001")
}
