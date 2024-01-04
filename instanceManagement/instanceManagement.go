package instanceManagement

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	game "Engee-Game/gamedummy"
	"Engee-Game/utils"

	"github.com/gorilla/websocket"
)

var instances map[string]GameInstance

func PrepareInstancing(port string) {
	instances = make(map[string]GameInstance)

	info := utils.StringPair{
		First:  "test",
		Second: "localhost:" + port,
	}

	body, err := json.Marshal(info)
	if err != nil {
		log.Fatalf("Could not register game mode (body): %v", err)
	}

	reqBody := bytes.NewReader(body)

	request, err := http.NewRequest("POST", "http://localhost:8090/gameModes", reqBody)
	if err != nil {
		log.Fatalf("Could not register game mode (request): %v", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatalf("Could not register game mode (sent): %v", err)
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		log.Fatalf("Could not register game mode (code): %v", err)
	}
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

func AddPlayerToInstance(rid string, uid string, conn *websocket.Conn) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	err = instance.JoinPlayer(uid, conn)
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
