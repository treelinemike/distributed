// config.go
//
// Handle loading Raft cluster configuration from yaml file.
//
// Author: M. Kokko
// Updated: 24-Jan-2025

package raftcommon

import (
	"log"
	"net/rpc"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type (
	NetworkAddress struct {
		Address string `mapstructure:"Address"`
		Port    string `mapstructure:"Port"`
		Handle  *rpc.Client
	}

	RaftConfig struct {
		Servers             map[string]NetworkAddress `mapstructure:"Servers"`
		Timeout             Timeout                   `mapstructure:"Timeout"`
		NVStateFilenameBase string                    `mapstructure:"NVStateFilenameBase"`
	}
)

// see: https://betterprogramming.pub/parsing-and-creating-yaml-in-go-crash-course-2ec10b7db850
func LoadRaftConfig(configFileName string, Servers map[string]NetworkAddress, t *Timeout, nvfilebase *string) error {

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
	*nvfilebase = cfg.NVStateFilenameBase

	// done
	return nil
}
