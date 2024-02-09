package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const configPath = "./config.json"

type Config struct {
	GameServer     string `json:"game_server"`
	GamePort       string `json:"game_port"`
	RegistryServer string `json:"registry_server"`
	RegistryPort   string `json:"registry_port"`
	GameMode       string `json:"game_mode"`
}

func ReadConfig() Config {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Could not read config file on launch: %v", err))
	}

	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		panic(fmt.Sprintf("Could not parse config file on launch: %v", err))
	}

	return payload
}
