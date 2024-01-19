package gamedummy

import (
	"encoding/json"
	"fmt"
	"log"
)

func (game *GameDummy) HandleMessage(uid string, message []byte) {
	var msg DummyMessage
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Printf("[Error] could not unmarshal message: %v", err)
		return
	}

	err = routeMessage(msg, game)
	if err != nil {
		log.Printf("[Error] could not route message: %v", err)
		return
	}
}

func routeMessage(msg DummyMessage, game *GameDummy) error {
	switch msg.Type {
	case "Connect":
		return game.Connect(msg.Content)
	case "Control":
		return game.Control(msg.Content)
	case "Test":
		return game.Test(msg.Content)
	default:
		return fmt.Errorf("invalid dummy message type: %v", msg.Type)
	}
}
