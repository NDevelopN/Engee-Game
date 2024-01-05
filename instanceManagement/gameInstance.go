package instanceManagement

type GameInstance interface {
	SetRules(rules string) error
	EndGame() error
	StartGame() error
	PauseGame() error
	ResetGame() error
	JoinPlayer(uid string) error
	RemovePlayer(uid string) error
	AddListener(listener func([]byte) error) (string, error)
	RemoveListener(id string) error
	HandleMessage([]byte)
}
