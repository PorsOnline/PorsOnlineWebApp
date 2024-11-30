package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func readConfig(configPath string) (Config, error) {
	var c Config
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func MustReadConfig(configPath string) Config {
	c, err := readConfig(configPath)
	if err != nil {
		log.Fatal("Can not read config file.", err.Error())
	}
	return c
}
