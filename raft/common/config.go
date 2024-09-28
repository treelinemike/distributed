package common

import (
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type (
	NetworkAddress struct {
		Address string `mapstructure:"Address"`
		Port    string `mapstructure:"Port"`
	}
	ServerConfig struct {
		Servers map[string]NetworkAddress `mapstructure:"Servers"`
	}
)

// see: https://betterprogramming.pub/parsing-and-creating-yaml-in-go-crash-course-2ec10b7db850
func LoadRaftConfig(configFileName string, Servers map[string]NetworkAddress) error {

	configFile, err := os.ReadFile(configFileName)
	if err != nil {
		log.Fatal("Cannot read config file!")
	}

	var yamlIn interface{}

	err = yaml.Unmarshal(configFile, &yamlIn)
	if err != nil {
		log.Fatal(err)
	}

	var cfg ServerConfig

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &cfg})
	err = decoder.Decode(yamlIn)
	if err != nil {
		log.Fatal(err)
	}
	for key := range cfg.Servers {
		Servers[key] = cfg.Servers[key]
	}

	return nil
}
