package instanceManagement

import "github.com/gorilla/websocket"

type GameInstance interface {
	SetRules(rules string) error
	EndGame() error
	StartGame() error
	PauseGame() error
	ResetGame() error
	JoinPlayer(uid string, conn *websocket.Conn) error
	RemovePlayer(uid string) error
}
