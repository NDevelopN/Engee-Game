package instanceManagement

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	game "Engee-Game/gamedummy"
	"Engee-Game/utils"
)

var instances map[string]GameInstance

const gameMode = "test"
const serverAddr = "http://localhost:8090"

const HeartbeatInterval = 3 * time.Second

func PrepareInstancing(gameAddr string) {
	instances = make(map[string]GameInstance)

	info := utils.StringPair{
		First:  gameMode,
		Second: gameAddr,
	}

	body, err := json.Marshal(info)
	if err != nil {
		log.Fatalf("Could not register game mode (body): %v", err)
	}

	reqBody := bytes.NewReader(body)

	request, err := http.NewRequest("POST", serverAddr+"/gameModes", reqBody)
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

	hbRequest, err := http.NewRequest("POST", serverAddr+"/gameModes/"+gameMode, nil)
	if err != nil {
		log.Fatalf("Could not prepare heartbeat message: %v", err)
	}

	go SendHeartbeats(hbRequest)
}

func SendHeartbeats(request *http.Request) error {
	for {
		_, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalf("[Error] failed to send heartbeat message: %v", err)
		}

		time.Sleep(HeartbeatInterval)
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

	instance, err := game.CreateDefaultGame(rid)
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

func AddPlayerToInstance(rid string, uid string, listener func(message []byte) error) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	err = instance.AddPlayer(uid, listener)
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

func MessageHandleInstance(rid string, uid string, message []byte) {
	instance, err := getInstance(rid)
	if err != nil {
		log.Printf("[Error] could not get instnace to handle message: %v", err)
		return
	}

	instance.HandleMessage(uid, message)
}

func getInstance(rid string) (GameInstance, error) {
	instance, found := instances[rid]
	if !found {
		return instance, fmt.Errorf("game does not exist for room %s", rid)
	}

	return instance, nil
}
