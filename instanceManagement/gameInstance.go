package instanceManagement

type GameInstance interface {
	SetRules(rules string) error
	EndGame() error
	StartGame() error
	PauseGame() error
	ResetGame() error
	JoinPlayer(uid string) error
	RemovePlayer(uid string) error
}
