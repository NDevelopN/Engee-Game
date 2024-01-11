package gamedummy

import (
	"encoding/json"
	"fmt"
	"log"
)

type DummyMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (game *GameDummy) SendRulesUpdate() error {
	content, err := json.Marshal(game.Rules)
	if err != nil {
		return err
	}

	msg := DummyMessage{
		Type:    "Rules",
		Content: string(content),
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = sendMessage(game.Listeners, message)
	if err != nil {
		log.Printf("[Error] Failed to send Rules update: %v", err)
	}

	return nil
}

func (game *GameDummy) SendStatusUpdate() error {
	msg := DummyMessage{
		Type:    "Status",
		Content: fmt.Sprint(game.Status),
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = sendMessage(game.Listeners, message)
	if err != nil {
		log.Printf("[Error] Failed to send Rules update: %v", err)
	}

	return nil
}

func (game *GameDummy) SendPlayerUpdate() error {
	players, err := json.Marshal(game.Players)
	if err != nil {
		return err
	}

	msg := DummyMessage{
		Type:    "Players",
		Content: string(players),
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = sendMessage(game.Listeners, message)
	if err != nil {
		log.Printf("[Error] Failed to send Rules update: %v", err)
	}

	return nil
}

func (game *GameDummy) SendXUpdate(x string, update string) error {
	msg := DummyMessage{
		Type:    x,
		Content: update,
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = sendMessage(game.Listeners, message)
	if err != nil {
		log.Printf("[Error] Failed to send generic update: %v", err)
	}

	return nil
}

func sendMessage(listeners map[string]func([]byte) error, message []byte) error {
	var err error

	if len(listeners) == 0 {
		log.Printf("[Alert] No listeners registered for game")
		return nil
	}

	fail := false

	for _, listener := range listeners {
		err = listener(message)
		if err != nil {
			fail = true
			log.Printf("[Error] Sending message to listener: %s", err)
		}
	}

	if fail {
		return fmt.Errorf("failed to send message to all listeners")
	}

	return nil
}
