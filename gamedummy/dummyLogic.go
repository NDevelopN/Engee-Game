package gamedummy

import (
	"fmt"

	sErr "Engee-Game/stockErrors"
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
	Listeners map[string]func([]byte) error
}

func CreateDefaultGame(rid string) (*GameDummy, error) {
	nGame := new(GameDummy)
	nGame.RID = rid
	nGame.Rules = dummyRules
	nGame.Status = NEW
	nGame.Players = []string{}
	nGame.Listeners = make(map[string]func([]byte) error)

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

func (dummy *GameDummy) AddPlayer(uid string, listener func([]byte) error) error {
	err := checkValidGame(dummy, []int{})
	if err != nil {
		return err
	}

	for _, plr := range dummy.Players {
		if plr == uid {
			return &sErr.MatchFoundError[string]{
				Space: "Game Players",
				Field: "UID",
				Value: plr,
			}
		}
	}

	dummy.Players = append(dummy.Players, uid)
	dummy.Listeners[uid] = listener

	return dummy.SendStatusUpdate()
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
			delete(dummy.Listeners, uid)
			return nil
		}
	}

	return &sErr.MatchNotFoundError[string]{
		Space: "Game Players",
		Field: "UID",
		Value: uid,
	}
}

func (dummy *GameDummy) EndGame() error {
	err := checkValidGame(dummy, []int{})
	if err != nil {
		return err
	}

	dummy.Listeners = nil

	return nil
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

	return sErr.InvalidValueError[string]{
		Field: "Command",
		Value: message,
	}
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
		return &sErr.EmptyValueError{
			Field: "Game Rules",
		}
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
			return &sErr.InvalidValueError[int]{
				Field: "Status",
				Value: dummy.Status,
			}
		}
	}

	return nil
}
