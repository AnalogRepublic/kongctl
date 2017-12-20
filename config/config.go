package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

type Context struct {
	Host string `yaml:"host"`
}

type Config struct {
	*viper.Viper
	FileData *ConfigFileData
}

type ConfigFileData struct {
	CurrentContext string             `yaml:"current_context"`
	Contexts       map[string]Context `yaml:"contexts"`
}

var config *Config

const (
	appName = "kongctl"
)

func GetConfig() *Config {
	return config
}

func Init() error {
	config = &Config{
		Viper: viper.New(),
	}

	config.setup()

	err := config.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		return errors.Wrap(err, "Fatal error loading config file")
	}

	fmt.Println("Using config from here: " + config.ConfigFileUsed())

	config.setDefaults()

	// This unmarshals the config file so we can
	// manipulate it and re-save it.
	err = config.loadFile()

	if err != nil { // Handle errors reading the config file
		return errors.Wrap(err, "Fatal error loading config file")
	}

	return nil
}

func (conf *Config) GetCurrentContext() (Context, error) {
	context := Context{}

	key := conf.FileData.CurrentContext
	context, exists := conf.FileData.Contexts[key]

	if !exists {
		return context, errors.New("You have not defined the current context in your config file")
	}

	return context, nil
}

func (conf *Config) setup() {
	// Load in YAML files
	conf.SetConfigType("yaml")

	// Allow us to set config using Environment Variables
	conf.SetEnvPrefix("")
	conf.AutomaticEnv()
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Want to make sure we're loading in the config file at
	// the right paths.
	conf.SetConfigName(fmt.Sprintf("%s", appName)) // name of config file (without extension)

	// Want to read in config from the home dir
	conf.AddConfigPath(fmt.Sprintf("$HOME/.%s", appName))

	// Always look for a config file in the current dir
	conf.AddConfigPath(".")
}

func (conf *Config) setDefaults() {
	// By default we'll assume that we're running the Kong
	// service on the same host machine.
	conf.SetDefault("host", "http://127.0.0.1:8001")
}

func (conf *Config) loadFile() error {
	file := config.ConfigFileUsed()
	data, err := ioutil.ReadFile(file)

	fileData := &ConfigFileData{}

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, fileData)

	if err != nil {
		return err
	}

	conf.FileData = fileData

	return nil
}

func (conf *Config) SaveFile() error {
	file := config.ConfigFileUsed()

	output, err := yaml.Marshal(config.FileData)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, output, 0644)
}
