package server

import (
	"Engee-Game/instanceManagement"
	"Engee-Game/utils"
	"strings"

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
	rid := ids[0]
	uid := ids[1]

	conn, err := upgradeConnection(w, r)
	if err != nil {
		http.Error(w, "Failed to upgrade to websocket connection", http.StatusInternalServerError)
		log.Printf("[Error] Upgrading connection: %v", err)
		return
	}

	conn.SetCloseHandler(handleClose)

	go listenWhileConnected(rid, uid, conn)

	err = instanceManagement.AddPlayerToInstance(ids[0], ids[1],
		(func(message []byte) error {
			err = conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return fmt.Errorf("could not send message on ws connection: %v", err)
			}
			return nil
		}))

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

func listenWhileConnected(rid string, uid string, conn *websocket.Conn) {
	acceptInput(rid, uid, conn)

	err := instanceManagement.RemovePlayerFromInstance(rid, uid)
	if err != nil {
		log.Printf("[Error] Removing player after conneciton ended: %v", err)
	}
}

const badMsgThreshold = 10

func acceptInput(rid string, uid string, conn *websocket.Conn) {
	badMsgCount := 0

	for {
		mType, data, err := conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") ||
				strings.Contains(err.Error(), "connection closed") {

				log.Printf("[Error] Network connection closed:  %v", err)
				return
			}

			log.Printf("[Error] reading input from player %s: %v", uid, err)
			badMsgCount++
			continue
		}

		if mType != websocket.TextMessage {
			log.Printf("[Error] message type not supported: %v", mType)
			badMsgCount++
			continue
		}

		badMsgCount = 0

		instanceManagement.MessageHandleInstance(rid, uid, data)
	}
}
