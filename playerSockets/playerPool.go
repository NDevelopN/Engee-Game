package playerSockets

import (
	im "Engee-Game/instanceManagement"
	"strings"

	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type SockPool struct {
	mutex       sync.Mutex
	connections map[string]*websocket.Conn
	listeners   map[string]string
}

var gamePools map[string]*SockPool = make(map[string]*SockPool)

func Instantiate(rid string) error {
	_, found := gamePools[rid]
	if found {
		return fmt.Errorf("socket Pool for that game already exists")
	}

	newPool := new(SockPool)
	newPool.mutex = sync.Mutex{}
	newPool.connections = map[string]*websocket.Conn{}
	newPool.listeners = map[string]string{}

	gamePools[rid] = newPool

	return nil
}

func CloseAll(rid string) error {
	pool, found := gamePools[rid]
	if !found {
		return fmt.Errorf("no socket pool to close")
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for _, conn := range pool.connections {
		err := conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func Delete(rid string) error {
	err := CloseAll(rid)
	if err != nil {
		return err
	}

	delete(gamePools, rid)

	return nil
}

func AddPlayerToPool(rid string, uid string, conn *websocket.Conn) error {
	log.Printf("Adding player %s to pool %s", uid, rid)
	pool, found := gamePools[rid]
	if !found {
		err := Instantiate(rid)
		if err != nil {
			return err
		}
		pool = gamePools[rid]
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	_, found = pool.connections[uid]
	if found {
		return fmt.Errorf("connection already established for player %q", uid)
	}

	lid, err := im.AddListenerToInstance(rid, (func(message []byte) error {
		return SendTo(rid, message, uid)
	}))

	if err != nil {
		return err
	}

	pool.listeners[uid] = lid
	pool.connections[uid] = conn
	gamePools[rid] = pool

	go AcceptInput(rid, uid, conn)

	return nil
}

func AcceptInput(rid string, uid string, conn *websocket.Conn) {
	for {
		mType, data, err := conn.ReadMessage()

		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("[Error] Network connection closed:  %v", err)
				return
			}
			log.Printf("[Error] reading input from player %s: %v", uid, err)
			continue
		}

		if mType != websocket.TextMessage {
			log.Printf("[Error] message type not supported: %v", mType)
			continue
		}

		im.MessageHandleInstance(rid, data)
	}
}

func RemovePlayerFromPool(rid string, uid string) error {
	pool, found := gamePools[rid]
	if !found {
		return fmt.Errorf("no socket pool exists for that game")
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	_, found = pool.connections[uid]
	if found {
		return fmt.Errorf("could not find connection for user %q", uid)
	}

	err := im.RemoveListenerFromInstance(rid, pool.listeners[uid])
	if err != nil {
		return err
	}

	delete(pool.connections, uid)
	gamePools[rid] = pool

	return nil
}
