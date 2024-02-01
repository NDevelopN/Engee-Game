package main

import (
	"Engee-Game/config"
	"Engee-Game/instanceManagement"
	"Engee-Game/server"
)

func main() {
	config := config.ReadConfig()
	instanceManagement.PrepareInstancing(config)
	server.Serve(config.GameServer.Port)
}
