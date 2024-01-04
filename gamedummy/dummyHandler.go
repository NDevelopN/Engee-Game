package gamedummy

import (
	"encoding/json"
	"fmt"
	"log"
)

func Handle(messageType int, data []byte, err error, game *GameDummy) {
	for {
		if err != nil {
			log.Printf("[CLOSE] Connection closed: %v", err)
			return
		}

		if messageType != 1 { //websocket.TextMessage
			log.Printf("[Error] Received unexpected message type: %v", messageType)
			continue
		}

		var msg = DummyMessage{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Printf("[Error] Cannot unmarshal received message: %v", err)
			continue
		}

		err = RouteMessage(msg, game)
		if err != nil {
			log.Printf("[Error] Handling message: %v", err)
			continue
		}
	}
}

func RouteMessage(msg DummyMessage, game *GameDummy) error {
	log.Printf("Routing message: %v", msg)
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
