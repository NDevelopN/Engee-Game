package main

import (
	"Engee-Game/instanceManagement"
	"Engee-Game/server"
)

const port = "8091"
const address = "http://localhost:" + port

func main() {
	instanceManagement.PrepareInstancing(address)
	server.Serve(port)
}
