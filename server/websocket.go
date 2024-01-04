package server

import (
	"Engee-Game/instanceManagement"
	pSock "Engee-Game/playerSockets"
	"Engee-Game/utils"

	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func JoinPlayer(c *gin.Context) {
	w := c.Writer
	r := c.Request

	ids := utils.GetRequestIDs(c.Request)

	conn, err := upgradeConnection(w, r)
	if err != nil {
		http.Error(w, "Failed to upgrade to websocket connection", http.StatusInternalServerError)
		log.Printf("[Error] Upgrading connection: %v", err)
		return
	}

	conn.SetCloseHandler(handleClose)

	err = pSock.AddPlayerToPool(ids[0], ids[1], conn)
	if err != nil {
		http.Error(w, "Failed to add player to pool", http.StatusInternalServerError)
		log.Printf("[Error] Adding player to connection pool: %v", err)
		conn.Close()
		return
	}

	err = instanceManagement.AddPlayerToInstance(ids[0], ids[1], conn)
	if err != nil {
		http.Error(w, "Failed to add player to game", http.StatusInternalServerError)
		log.Printf("[Error] Adding player to game: %v", err)
		conn.Close()
		return
	}
}

func upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{}
	upgrader.CheckOrigin =
		func(r *http.Request) bool {
			return true
		}

	return upgrader.Upgrade(w, r, nil)
}

func handleClose(code int, text string) error {
	if code == websocket.CloseNoStatusReceived {
		text = "without status"
	}

	return fmt.Errorf("connection closed: %s", text)
}
