package gamedummy

import "fmt"

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
}

func CreateDefaultGame() *GameDummy {
	nGame := new(GameDummy)
	nGame.Rules = dummyRules
	nGame.Status = NEW
	nGame.Players = []string{}

	return nGame
}

func (dummy *GameDummy) SetRules(rules string) error {
	err := checkValidGame(dummy, []int{NEW, RESET})
	if err != nil {
		return err
	}

	dummy.Rules = rules

	return nil

}

func (dummy *GameDummy) StartGame() error {
	err := checkValidGame(dummy, []int{NEW, RESET})
	if err != nil {
		return err
	}

	dummy.Status = ACTIVE

	return nil
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

	return nil
}

func (dummy *GameDummy) ResetGame() error {
	err := checkValidGame(dummy, []int{ACTIVE, PAUSED})
	if err != nil {
		return err
	}

	dummy.Status = RESET

	return nil
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
			return nil
		}
	}

	return fmt.Errorf("player not found in the game")
}

func (dummy *GameDummy) EndGame() error {
	err := checkValidGame(dummy, []int{})
	if err != nil {
		return err
	}

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
