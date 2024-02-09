package instanceManagement

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"Engee-Game/config"
	game "Engee-Game/gamedummy"
	sErr "Engee-Game/stockErrors"
	"Engee-Game/utils"
)

var instances map[string]GameInstance

const HeartbeatInterval = 3 * time.Second

func PrepareInstancing(config config.Config) {
	instances = make(map[string]GameInstance)

	gameAddr := fmt.Sprintf("http://%s:%s", config.GameServer, config.GamePort)
	regAddr := fmt.Sprintf("http://%s:%s", config.RegistryServer, config.RegistryPort)
	info := utils.StringPair{
		First:  config.GameMode,
		Second: gameAddr,
	}

	body, err := json.Marshal(info)
	if err != nil {
		log.Fatalf("Could not register game mode (body): %v", err)
	}

	reqBody := bytes.NewReader(body)

	request, err := http.NewRequest("POST", regAddr+"/gameModes", reqBody)
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

	hbRequest, err := http.NewRequest("POST", regAddr+"/gameModes/"+config.GameMode, nil)
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
		return &sErr.EmptyValueError{
			Field: "RID",
		}
	}

	_, found := instances[rid]
	if found {
		return &sErr.MatchFoundError[string]{
			Space: "Games",
			Field: "RID",
			Value: rid,
		}
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
		return instance, &sErr.MatchNotFoundError[string]{
			Space: "Games",
			Field: "RID",
			Value: rid,
		}
	}

	return instance, nil
}
