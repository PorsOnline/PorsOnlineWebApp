package config

import (
	"encoding/json"
	"log"
	"os"
)

func ReadConfig(configPath string) (Config, error) {
	var c Config
	all, err := os.ReadFile(configPath)
	if err != nil {
		return c, err
	}

	return c, json.Unmarshal(all, &c)
}

func MustReadConfig(configPath string) Config {
	c, err := ReadConfig(configPath)
	if err != nil {
		log.Fatal("Can not read config file.", err.Error())
	}
	return c
}

// func readConfig(configPath string) (Config, error) {
// 	var c Config
// 	yamlFile, err := os.ReadFile(configPath)
// 	if err != nil {
// 		return c, err
// 	}
// 	err = yaml.Unmarshal(yamlFile, &c)
// 	if err != nil {
// 		return c, err
// 	}
// 	return c, nil
// }
