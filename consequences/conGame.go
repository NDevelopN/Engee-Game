package consequences

import (
	"encoding/json"
	"fmt"

	sErr "Engee-Game/stockErrors"
)

var defaultPrompts = []string{
	"Person 1",
	"Person 2",
	"Location",
	"Action 1",
	"Action 2",
	"Consequences",
}

type ConRules struct {
	ShuffleType int
	Prompts     []string
}

type ConGame struct {
	RID          string
	Status       int
	Paused       bool
	Rules        ConRules
	PlayerStatus map[string]bool
	ReplySet     map[string][]string
	Shuffled     map[string][]string
	Listeners    map[string]func([]byte) error
}

var defaultRules = ConRules{
	ShuffleType: DEFAULT_SHUFFLE,
	Prompts:     defaultPrompts,
}

const (
	WAITING int = 0
	PROMPTS int = 1
	REPLIES int = 2
)

const (
	DEFAULT_SHUFFLE int = 0
)

func CreateDefaultGame(rid string) (*ConGame, error) {
	nGame := new(ConGame)
	nGame.RID = rid
	nGame.Status = WAITING
	nGame.Paused = false
	nGame.Rules = defaultRules
	nGame.PlayerStatus = make(map[string]bool)
	nGame.ReplySet = make(map[string][]string)
	nGame.Shuffled = nil
	nGame.Listeners = make(map[string]func([]byte) error)

	return nGame, nil
}

func (game *ConGame) SetRules(rString string) error {
	var rules ConRules
	err := json.Unmarshal([]byte(rString), &rules)
	if err != nil {
		return err
	}

	err = game.ResetGame()
	if err != nil {
		return fmt.Errorf("could not reset game to apply rules update: %w", err)
	}

	game.Rules = rules

	return nil
}

//TODO can this just be removed now?
func (game *ConGame) StartGame() error {
	err := game.SendPrompts()
	if err != nil {
		return err
	}

	game.Status = PROMPTS

	return game.SendStatus()
}

func (game *ConGame) PauseGame() error {
	game.Paused = !game.Paused

	return game.SendStatus()
}

func (game *ConGame) ResetGame() error {
	game.Status = WAITING
	game.PlayerStatus = make(map[string]bool)
	game.ReplySet = make(map[string][]string)
	game.Shuffled = nil

	return game.SendStatus()
}

func (game *ConGame) EndGame() error {

	//TODO send update to all players and kill WS

	return nil
}

func (game *ConGame) AddPlayer(uid string, listener func([]byte) error) error {
	_, found := game.ReplySet[uid]
	if found {
		return &sErr.MatchFoundError[string]{
			Space: "Game Players",
			Field: "UID",
			Value: uid,
		}
	}

	game.ReplySet[uid] = make([]string, len(game.Rules.Prompts))

	//TODO REMOVE
	if game.Status == WAITING {
		err := game.PlayerReady(uid)
		if err != nil {
			return err
		}
	}

	game.Listeners[uid] = listener

	err := game.SendPromptsTo(uid)
	if err != nil {
		return err
	}
	return game.SendStatusTo(uid)
}

func (game *ConGame) RemovePlayer(uid string) error {
	_, found := game.ReplySet[uid]
	if !found {
		return &sErr.MatchNotFoundError[string]{
			Space: "Game Players",
			Field: "UID",
			Value: uid,
		}
	}

	delete(game.ReplySet, uid)
	delete(game.Shuffled, uid)

	_, found = game.Listeners[uid]
	if !found {
		return &sErr.MatchNotFoundError[string]{
			Space: "Player Listeners",
			Field: "UID",
			Value: uid,
		}
	}

	delete(game.Listeners, uid)

	return nil
}

func (game *ConGame) GetStatus(uid string) error {
	return game.SendStatusTo(uid)
}

func (game *ConGame) GetPrompts(uid string) error {
	return game.SendPromptsTo(uid)
}

func (game *ConGame) AcceptReply(uid string, replies []string) error {
	_, found := game.ReplySet[uid]
	if !found {
		return &sErr.MatchNotFoundError[string]{
			Space: "Player Listeners",
			Field: "UID",
			Value: uid,
		}
	}

	if len(replies) != len(game.Rules.Prompts) {
		return &sErr.InvalidSetLengthError{
			Set:      "Replies",
			Expected: len(game.Rules.Prompts),
			Actual:   len(replies),
		}
	}

	for index, reply := range replies {
		if reply == "" {
			return &sErr.EmptyValueError{
				Field: fmt.Sprintf("Reply [%d]", index),
			}
		}
	}

	game.ReplySet[uid] = replies

	ready := 0
	for _, replies = range game.ReplySet {
		if len(replies) == len(game.Rules.Prompts) {
			ready++
		}
	}

	if ready == len(game.ReplySet) {
		err := game.switchToShuffle()
		if err != nil {
			return err
		}
	}

	return nil
}

func (game *ConGame) PlayerReady(uid string) error {
	if game.Status != WAITING {
		return &sErr.InvalidValueError[int]{
			Field: "Game Status",
			Value: game.Status,
		}
	}

	game.PlayerStatus[uid] = true

	for _, status := range game.PlayerStatus {
		if !status {
			return nil
		}
	}

	game.Status = PROMPTS
	return game.SendStatus()
}

func (game *ConGame) switchToShuffle() error {
	err := game.shuffleReplies()
	if err != nil {
		return err
	}

	game.Status = REPLIES

	for uid := range game.Shuffled {
		err = game.SendShuffledTo(uid)
		if err != nil {
			return err
		}
	}

	err = game.SendStatus()
	if err != nil {
		return err
	}

	game.Status = WAITING
	return nil
}

func (game *ConGame) GetShuffledReplies(uid string) error {
	_, found := game.ReplySet[uid]
	if !found {
		return &sErr.MatchNotFoundError[string]{
			Space: "Game Players",
			Field: "UID",
			Value: uid,
		}
	}

	if game.Shuffled == nil {
		return &sErr.EmptyValueError{
			Field: "Shuffled Replies",
		}
	}

	return game.SendShuffledTo(uid)
}

func (game *ConGame) shuffleReplies() error {
	shuffled := make(map[string][]string)

	var intToUID []string
	for uid, userSet := range game.ReplySet {
		if userSet == nil {
			return &sErr.MatchNotFoundError[string]{
				Space: "Reply Sets",
				Field: "UID",
				Value: uid,
			}
		}
		intToUID = append(intToUID, uid)
	}

	for uIndex, uid := range intToUID {
		pIndex := 0
		shuffleSet := make([]string, len(game.Rules.Prompts))
		for pIndex < len(game.Rules.Prompts) {
			shuffleIndex := (uIndex + pIndex) % len(game.ReplySet)
			shuffleUID := intToUID[shuffleIndex]
			nextReply := game.ReplySet[shuffleUID][pIndex]
			if nextReply == "" {
				return sErr.EmptyValueError{
					Field: fmt.Sprintf("Reply [%d,%d]", uIndex, pIndex),
				}
			}
			shuffleSet[pIndex] = nextReply
			pIndex++
		}

		shuffled[uid] = shuffleSet
	}

	game.Shuffled = shuffled

	return nil
}
