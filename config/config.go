package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const configPath = "./config.json"

type Address struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type Config struct {
	GameServer     Address `json:"game_server"`
	RegistryServer Address `json:"registry_server"`
	GameMode       string  `json:"game_mode"`
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
