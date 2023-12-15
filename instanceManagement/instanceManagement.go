package instanceManagement

import (
	"fmt"

	game "Engee-Game/gamedummy"
	"Engee-Game/websocket"
)

var instances map[string]GameInstance

func PrepareInstancing() {
	instances = make(map[string]GameInstance)
}

func CreateNewInstance(rid string) error {
	if rid == "" {
		return fmt.Errorf("empty RID provided")
	}

	_, found := instances[rid]
	if found {
		return fmt.Errorf("game already exists for room %s", rid)
	}

	instance := game.CreateDefaultGame()

	err := websocket.Instantiate(rid)
	if err != nil {
		return err
	}

	instances[rid] = instance

	return nil
}

func DeleteInstance(rid string) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	err = instance.EndGame()
	if err != nil {
		return err
	}

	err = websocket.Remove(rid)
	if err != nil {
		return fmt.Errorf("could not close connection pool: %v", err)
	}
	delete(instances, rid)

	return nil
}

func StartInstance(rid string) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	err = instance.StartGame()
	if err != nil {
		return err
	}

	instances[rid] = instance

	return nil
}

func SetInstanceRules(rid string, rules string) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	err = instance.SetRules(rules)
	if err != nil {
		return err
	}

	instances[rid] = instance
	return nil
}

func PauseInstance(rid string) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	err = instance.PauseGame()
	if err != nil {
		return err
	}

	instances[rid] = instance

	return nil
}

func ResetInstance(rid string) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	err = instance.ResetGame()
	if err != nil {
		return err
	}

	instances[rid] = instance
	return nil
}

func RemovePlayerFromInstance(rid string, uid string) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	err = instance.RemovePlayer(uid)
	if err != nil {
		return err
	}

	websocket.RemovePlayerFromPool(rid, uid)

	instances[rid] = instance
	return nil
}

func getInstance(rid string) (GameInstance, error) {
	instance, found := instances[rid]
	if !found {
		return instance, fmt.Errorf("game does not exist for room %s", rid)
	}

	return instance, nil
}
