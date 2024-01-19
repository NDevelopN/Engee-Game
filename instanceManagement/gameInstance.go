package instanceManagement

type GameInstance interface {
	SetRules(rules string) error
	EndGame() error
	StartGame() error
	PauseGame() error
	ResetGame() error
	AddPlayer(id string, listener func([]byte) error) error
	RemovePlayer(uid string) error
	HandleMessage(uid string, message []byte)
}
