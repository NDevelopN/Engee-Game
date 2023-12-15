package instanceManagement

type GameInstance interface {
	SetRules(rules string) error
	EndGame() error
	StartGame() error
	PauseGame() error
	ResetGame() error
	RemovePlayer(uid string) error
}
