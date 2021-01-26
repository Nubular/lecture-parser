package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Config struct for the module, update as necessary
type Config struct {
	CacheFiles bool `json:"cacheFiles"`
}

func loadConfig() (*Config, error) {
	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configFile, err := os.Open(filepath.Join(absPath, ".config/config.json"))
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	config := Config{}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil

}