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
	Timeout struct {
		Min_ms int `mapstructure:"Min_ms"`
		Max_ms int `mapstructure:"Max_ms"`
	}
	RaftConfig struct {
		Servers         map[string]NetworkAddress `mapstructure:"Servers"`
		Timeout         Timeout                   `mapstructure:"Timeout"`
		NVStateFilename string                    `mapstructure:"NVStateFilename"`
	}
)

// see: https://betterprogramming.pub/parsing-and-creating-yaml-in-go-crash-course-2ec10b7db850
func LoadRaftConfig(configFileName string, Servers map[string]NetworkAddress, t *Timeout, nvstatefilename *string) error {

	configFile, err := os.ReadFile(configFileName)
	if err != nil {
		log.Fatal("Cannot read config file!")
	}

	var yamlIn interface{}

	err = yaml.Unmarshal(configFile, &yamlIn)
	if err != nil {
		log.Fatal(err)
	}

	// decode config file into mapstructure
	var cfg RaftConfig
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &cfg})
	err = decoder.Decode(yamlIn)
	if err != nil {
		log.Fatal(err)
	}

	// store server list
	for key := range cfg.Servers {
		Servers[key] = cfg.Servers[key]
	}

	// store timeout
	t.Min_ms = cfg.Timeout.Min_ms
	t.Max_ms = cfg.Timeout.Max_ms

	// store NV state filename
	*nvstatefilename = cfg.NVStateFilename

	// done
	return nil
}
