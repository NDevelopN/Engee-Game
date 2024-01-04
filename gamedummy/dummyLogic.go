package gamedummy

import (
	pSock "Engee-Game/playerSockets"
	"fmt"
)

const dummyRules = "Rules"

const (
	NEW    int = 0
	ACTIVE int = 1
	PAUSED int = 2
	RESET  int = 3
	END    int = 4
)

type GameDummy struct {
	RID       string
	Rules     string
	Status    int
	Players   []string
	Listeners []func([]byte) error
}

func CreateDefaultGame(rid string) (*GameDummy, error) {
	nGame := new(GameDummy)
	nGame.RID = rid
	nGame.Rules = dummyRules
	nGame.Status = NEW
	nGame.Players = []string{}

	err := pSock.Instantiate(rid,
		func(messageType int, data []byte, err error) {
			Handle(messageType, data, err, nGame)
		})

	if err != nil {
		return nil, err
	}

	return nGame, nil
}

func (dummy *GameDummy) SetRules(rules string) error {
	err := checkValidGame(dummy, []int{NEW, RESET})
	if err != nil {
		return err
	}

	dummy.Rules = rules

	return dummy.SendRulesUpdate()
}

func (dummy *GameDummy) StartGame() error {
	err := checkValidGame(dummy, []int{NEW, RESET})
	if err != nil {
		return err
	}

	dummy.Status = ACTIVE

	return dummy.SendStatusUpdate()
}

func (dummy *GameDummy) PauseGame() error {
	err := checkValidGame(dummy, []int{ACTIVE, PAUSED})
	if err != nil {
		return err
	}

	if dummy.Status == ACTIVE {
		dummy.Status = PAUSED
	} else {
		dummy.Status = ACTIVE
	}

	return dummy.SendStatusUpdate()
}

func (dummy *GameDummy) ResetGame() error {
	err := checkValidGame(dummy, []int{ACTIVE, PAUSED})
	if err != nil {
		return err
	}

	dummy.Status = RESET

	return dummy.SendStatusUpdate()
}

func (dummy *GameDummy) JoinPlayer(uid string) error {
	err := checkValidGame(dummy, []int{})
	if err != nil {
		return err
	}

	for _, plr := range dummy.Players {
		if plr == uid {
			return fmt.Errorf("player already in the game")
		}
	}

	dummy.Players = append(dummy.Players, uid)

	return dummy.SendPlayerUpdate()
}

func (dummy *GameDummy) AddListener(listener func([]byte) error) error {
	dummy.Listeners = append(dummy.Listeners, listener)

	return nil
}

func (dummy *GameDummy) RemovePlayer(uid string) error {
	err := checkValidGame(dummy, []int{})
	if err != nil {
		return err
	}

	players := dummy.Players
	end := len(players) - 1

	for index, plr := range players {
		if plr == uid {
			players[index] = players[end]
			dummy.Players = players[:end]
			pSock.RemovePlayerFromPool(dummy.RID, uid)
			return dummy.SendPlayerUpdate()
		}
	}

	return fmt.Errorf("player not found in the game")
}

func (dummy *GameDummy) EndGame() error {
	err := checkValidGame(dummy, []int{})
	if err != nil {
		return err
	}

	return pSock.CloseAll(dummy.RID)
}

func (dummy *GameDummy) Connect(message string) error {
	err := checkValidGame(dummy, []int{})
	if err != nil {
		return err
	}

	fmt.Printf("%s Connected ", message)

	dummy.SendStatusUpdate()

	return nil
}

func (dummy *GameDummy) Control(message string) error {
	switch message {
	case "Start":
		return dummy.StartGame()
	case "Pause":
		return dummy.PauseGame()
	case "Reset":
		return dummy.ResetGame()
	case "End":
		return dummy.EndGame()
	}

	return fmt.Errorf("unrecognised command: %q", message)

}

func (dummy *GameDummy) Test(message string) error {
	err := checkValidGame(dummy, []int{})
	if err != nil {
		return err
	}

	fmt.Printf("Message received: %s", message)

	dummy.SendXUpdate("Test", message+" (reply)")

	return nil
}

func checkValidGame(dummy *GameDummy, status []int) error {
	if dummy.Rules == "" {
		return fmt.Errorf("rules are not set")
	}

	if len(status) > 0 {
		validStatus := false

		for _, s := range status {
			if dummy.Status == s {
				validStatus = true
				continue
			}
		}

		if !validStatus {
			return fmt.Errorf("game is not in a valid state")
		}
	}

	return nil
}
