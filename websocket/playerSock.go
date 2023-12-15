package websocket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type SockPool struct {
	mutex       sync.Mutex
	connections map[string]*websocket.Conn
}

var gameSockets = make(map[string]*SockPool)

func Instantiate(rid string) error {
	_, found := gameSockets[rid]
	if found {
		return fmt.Errorf("connection pool already exists for room %q", rid)
	}

	newPool := new(SockPool)
	newPool.mutex = sync.Mutex{}
	newPool.connections = map[string]*websocket.Conn{}

	gameSockets[rid] = newPool

	return nil
}

func Remove(rid string) error {
	_, found := gameSockets[rid]
	if !found {
		return fmt.Errorf("connection pool for room %q not found", rid)
	}

	gameSockets[rid].mutex.Lock()
	defer gameSockets[rid].mutex.Unlock()

	for _, conn := range gameSockets[rid].connections {
		err := conn.Close()
		if err != nil {
			return err
		}
	}

	delete(gameSockets, rid)

	return nil
}

func AddPlayerToPool(rid string, uid string, conn *websocket.Conn) error {
	_, found := gameSockets[rid]
	if !found {
		return fmt.Errorf("could not find a socket pool for game %q", rid)
	}

	gameSockets[rid].mutex.Lock()
	defer gameSockets[rid].mutex.Unlock()

	_, found = gameSockets[rid].connections[uid]
	if found {
		return fmt.Errorf("connection already established for player %q in game %q", uid, rid)
	}

	gameSockets[rid].connections[uid] = conn

	return nil
}

func RemovePlayerFromPool(rid string, uid string) error {
	_, found := gameSockets[rid]
	if !found {
		return fmt.Errorf("could not find a socket pool for game %q", rid)
	}

	gameSockets[rid].mutex.Lock()
	defer gameSockets[rid].mutex.Unlock()

	_, found = gameSockets[rid].connections[uid]
	if found {
		return fmt.Errorf("could not find connection for user %q in game %q", uid, rid)
	}

	delete(gameSockets[rid].connections, uid)

	return nil
}

func SendAll(message []byte, rid string) error {
	_, found := gameSockets[rid]
	if !found {
		return fmt.Errorf("could not find a socket pool for game %q", rid)
	}

	if len(gameSockets[rid].connections) == 0 {
		return fmt.Errorf("no connections in socket pool of game %q", rid)
	}

	gameSockets[rid].mutex.Lock()
	defer gameSockets[rid].mutex.Unlock()

	for _, conn := range gameSockets[rid].connections {
		conn.WriteMessage(websocket.TextMessage, message)
	}

	return nil
}

func SendTo(message []byte, rid string, uid string) error {
	_, found := gameSockets[rid]
	if !found {
		return fmt.Errorf("coudl not find socket pool for game")
	}

	gameSockets[rid].mutex.Lock()
	defer gameSockets[rid].mutex.Unlock()

	conn, found := gameSockets[rid].connections[uid]
	if found {
		return fmt.Errorf("could not find connection for user %q in game %q", uid, rid)
	}

	conn.WriteMessage(websocket.TextMessage, message)

	return nil
}
