package main

import (
	"Engee-Game/instanceManagement"
	"Engee-Game/server"
)

const port = "8091"

func main() {
	instanceManagement.PrepareInstancing()
	server.Serve(port)
}