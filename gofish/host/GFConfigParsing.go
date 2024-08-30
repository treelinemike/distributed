package main

import (
	"log"
	"os"
	"sort"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

// see: https://betterprogramming.pub/parsing-and-creating-yaml-in-go-crash-course-2ec10b7db850
type (
	NetworkAddress struct {
		Address string `mapstructure:"Address"`
		Port    string `mapstructure:"Port"`
	}
)

type GFConfig struct {
	Host    NetworkAddress            `mapstructure:"Host"`
	Players map[string]NetworkAddress `mapstructure:"Players"`
}

func LoadGFConfig(configFileName string, HostIP *NetworkAddress, PlayerIP *[]NetworkAddress) error {

	configFile, err := os.ReadFile(configFileName)
	if err != nil {
		log.Fatal("Cannot read config file!")
	}

	var cfg GFConfig
	var yamlIn interface{}

	err = yaml.Unmarshal(configFile, &yamlIn)
	if err != nil {
		log.Fatal(err)
	}

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &cfg})
	err = decoder.Decode(yamlIn)
	if err != nil {
		log.Fatal(err)
	}

	// extract the host IP
	HostIP.Address = cfg.Host.Address
	HostIP.Port = cfg.Host.Port

	// extract player IPs
	// need to sort b/c Golang randomizes maps (thanks Golang...)
	// sorting helpful for grading, but not otherwise required
	var allKeys []string
	for player := range cfg.Players {
		allKeys = append(allKeys, player)
	}
	sort.Strings(allKeys)
	for _, thisKey := range allKeys {
		*PlayerIP = append(*PlayerIP, cfg.Players[thisKey])
	}

	return nil
}
