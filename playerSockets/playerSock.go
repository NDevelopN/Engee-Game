package playerSockets

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type SockPool struct {
	mutex       sync.Mutex
	connections map[string]*websocket.Conn
	handler     func(int, []byte, error)
}

var gamePools map[string]*SockPool = make(map[string]*SockPool)

func Instantiate(rid string, handler func(int, []byte, error)) error {
	_, found := gamePools[rid]
	if found {
		return fmt.Errorf("socket Pool for that game already exists")
	}

	newPool := new(SockPool)
	newPool.mutex = sync.Mutex{}
	newPool.connections = map[string]*websocket.Conn{}
	newPool.handler = handler

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
		return fmt.Errorf("no socket pool exists for that game")
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	_, found = pool.connections[uid]
	if found {
		return fmt.Errorf("connection already established for player %q", uid)
	}

	pool.connections[uid] = conn
	gamePools[rid] = pool

	go func() {
		for {
			pool.handler(pool.connections[uid].ReadMessage())
		}
	}()

	return nil
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

	delete(pool.connections, uid)
	gamePools[rid] = pool

	return nil
}

func SendAll(rid string, message []byte) error {
	pool, found := gamePools[rid]
	if !found {
		return fmt.Errorf("no socket pool exists for that game")
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if len(pool.connections) == 0 {
		return fmt.Errorf("no connections in socket pool")
	}

	for _, conn := range pool.connections {
		conn.WriteMessage(websocket.TextMessage, message)
	}

	return nil
}

func SendTo(rid string, message []byte, uid string) error {
	pool, found := gamePools[rid]
	if !found {
		return fmt.Errorf("no socket pool exists for that game")
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	conn, found := pool.connections[uid]
	if found {
		return fmt.Errorf("could not find connection for user %q", uid)
	}

	conn.WriteMessage(websocket.TextMessage, message)

	return nil
}
