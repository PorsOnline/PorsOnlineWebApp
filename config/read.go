package config

import (

	"encoding/json"
	"os"
	"log"
	"gopkg.in/yaml.v3"
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
		panic(err)


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


