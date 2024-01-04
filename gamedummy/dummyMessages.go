package gamedummy

import (
	pSock "Engee-Game/websocket"

	"encoding/json"
	"fmt"
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

	return pSock.SendAll(game.RID, message)
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

	return pSock.SendAll(game.RID, message)
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

	return pSock.SendAll(game.RID, message)
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

	return pSock.SendAll(game.RID, message)
}
