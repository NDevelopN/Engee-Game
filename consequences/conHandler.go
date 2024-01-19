package consequences

import (
	"encoding/json"
	"log"
)

func (game *ConGame) HandleMessage(uid string, message []byte) {
	var msg ConMessage
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Printf("[Error] could not unmarshal message: %v", err)
		return
	}

	switch msg.Type {
	case "Get":
		switch msg.Content {
		case "Status":
			err = game.GetStatus(uid)
			if err != nil {
				log.Printf("[Error] Could not get status update: %v", err)
				return
			}
		case "Prompts":
			err = game.GetPrompts(uid)
			if err != nil {
				log.Printf("[Error] Could not get prompts: %v", err)
				return
			}
		case "Shuffled":
			err = game.GetShuffledReplies(uid)
			if err != nil {
				log.Printf("[Error] Could not get shuffled replies: %v", err)
				return
			}
		default:
			log.Printf("[Error] Invalid information requested: %s", msg.Content)
			return
		}
	case "Reply":
		var replies []string
		err = json.Unmarshal([]byte(msg.Content), &replies)
		if err != nil {
			log.Printf("[Error] Could not parse reply: %v", err)
			return
		}

		err = game.AcceptReply(uid, replies)
		if err != nil {
			log.Printf("[Error] Could not accept reply: %v", err)
			return
		}
	case "Continue":
		err = game.PlayerReady(uid)
		if err != nil {
			log.Printf("[Error] Could not accept Continue: %v", err)
			return
		}

	default:
		log.Printf("[Error] Invalid message type: %s", msg.Type)
	}
}
