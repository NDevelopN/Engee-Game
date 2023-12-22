package gamedummy

import (
	pSock "Engee-Game/websocket"
	"fmt"

	"github.com/gorilla/websocket"
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
	Rules   string
	Status  int
	Players []string

	Pool *pSock.SockPool
}

func CreateDefaultGame() *GameDummy {
	nGame := new(GameDummy)
	nGame.Rules = dummyRules
	nGame.Status = NEW
	nGame.Players = []string{}

	nGame.Pool = pSock.Instantiate(
		func(conn *websocket.Conn) {
			Handle(conn, nGame)
		})

	return nGame
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

func (dummy *GameDummy) JoinPlayer(uid string, conn *websocket.Conn) error {
	err := checkValidGame(dummy, []int{})
	if err != nil {
		return err
	}

	for _, plr := range dummy.Players {
		if plr == uid {
			return fmt.Errorf("player already in the game")
		}
	}

	err = dummy.Pool.AddPlayerToPool(uid, conn)
	if err != nil {
		return err
	}

	return dummy.SendPlayerUpdate()
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
			dummy.Pool.RemovePlayerFromPool(uid)
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

	return dummy.Pool.CloseAll()
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

	return fmt.Errorf("Unrecognised command: %q", message)

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
