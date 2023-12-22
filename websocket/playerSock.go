package websocket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type SockPool struct {
	mutex       sync.Mutex
	connections map[string]*websocket.Conn
	handler     func(conn *websocket.Conn)
}

func Instantiate(handler func(conn *websocket.Conn)) *SockPool {
	newPool := new(SockPool)
	newPool.mutex = sync.Mutex{}
	newPool.connections = map[string]*websocket.Conn{}
	newPool.handler = handler

	return newPool
}

func (pool *SockPool) CloseAll() error {

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

func (pool *SockPool) AddPlayerToPool(uid string, conn *websocket.Conn) error {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	_, found := pool.connections[uid]
	if found {
		return fmt.Errorf("connection already established for player %q", uid)
	}

	pool.connections[uid] = conn

	go pool.handler(conn)

	return nil
}

func (pool *SockPool) RemovePlayerFromPool(uid string) error {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	_, found := pool.connections[uid]
	if found {
		return fmt.Errorf("could not find connection for user %q", uid)
	}

	delete(pool.connections, uid)

	return nil
}

func (pool *SockPool) SendAll(message []byte) error {
	if len(pool.connections) == 0 {
		return fmt.Errorf("no connections in socket pool")
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for _, conn := range pool.connections {
		conn.WriteMessage(websocket.TextMessage, message)
	}

	return nil
}

func (pool *SockPool) SendTo(message []byte, uid string) error {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	conn, found := pool.connections[uid]
	if found {
		return fmt.Errorf("could not find connection for user %q", uid)
	}

	conn.WriteMessage(websocket.TextMessage, message)

	return nil
}
