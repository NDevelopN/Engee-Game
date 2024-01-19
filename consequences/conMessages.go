package consequences

import (
	"encoding/json"
	"fmt"
	"log"
)

type ConMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (game *ConGame) SendStatus() error {
	message, err := game.packageStatusUpdate()
	if err != nil {
		return err
	}

	err = game.broadcastMessage(message)
	if err != nil {
		return fmt.Errorf("could not broadcast status update message: %v", err)
	}

	return nil
}

func (game *ConGame) SendStatusTo(uid string) error {
	message, err := game.packageStatusUpdate()
	if err != nil {
		return err
	}

	err = game.sendTo(uid, message)
	if err != nil {
		return fmt.Errorf("could not send status update message: %v", err)
	}

	return nil
}

func (game *ConGame) packageStatusUpdate() ([]byte, error) {
	msg := ConMessage{
		Type:    "Status",
		Content: fmt.Sprint(game.Status),
	}

	return json.Marshal(msg)
}

func (game *ConGame) SendPrompts() error {
	message, err := game.packagePrompts()
	if err != nil {
		return err
	}

	err = game.broadcastMessage(message)
	if err != nil {
		return fmt.Errorf("could not broadcast prompts message: %v", err)
	}

	return nil
}

func (game *ConGame) SendPromptsTo(uid string) error {
	message, err := game.packagePrompts()
	if err != nil {
		return err
	}

	err = game.sendTo(uid, message)
	if err != nil {
		return fmt.Errorf("could not send prompts message: %v", err)
	}

	return nil
}

func (game *ConGame) packagePrompts() ([]byte, error) {
	prompts, err := json.Marshal(game.Rules.Prompts)
	if err != nil {
		return nil, err
	}

	msg := ConMessage{
		Type:    "Prompts",
		Content: string(prompts),
	}

	return json.Marshal(msg)
}

func (game *ConGame) SendShuffledTo(uid string) error {
	message, err := game.packageShuffled(uid)
	if err != nil {
		return err
	}

	err = game.sendTo(uid, message)
	if err != nil {
		return fmt.Errorf("could not send shuffled replies: %v", err)
	}

	return nil
}

func (game *ConGame) packageShuffled(uid string) ([]byte, error) {
	shuffled, err := json.Marshal(game.Shuffled[uid])
	if err != nil {
		return nil, err
	}

	msg := ConMessage{
		Type:    "Shuffled",
		Content: string(shuffled),
	}

	return json.Marshal(msg)
}

func (game *ConGame) broadcastMessage(message []byte) error {
	if len(game.Listeners) == 0 {
		log.Printf("[Error] No listeners registered for game")
		return nil
	}

	fail := false

	for _, listener := range game.Listeners {
		err := listener(message)
		if err != nil {
			fail = true
			log.Printf("[Error] Sending message to listener: %v", err)
		}
	}

	if fail {
		return fmt.Errorf("failed to send message to all listeners")
	}

	return nil
}

func (game *ConGame) sendTo(uid string, message []byte) error {
	listener, found := game.Listeners[uid]
	if !found {
		return fmt.Errorf("no listener for found for uid: %s", uid)
	}

	return listener(message)
}
