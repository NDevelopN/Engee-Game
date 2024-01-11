package playerSockets

import (
	"fmt"

	"github.com/gorilla/websocket"
)

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
	if !found {
		return fmt.Errorf("could not find connection for user %q", uid)
	}

	conn.WriteMessage(websocket.TextMessage, message)

	return nil
}
