package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	im "Engee-Game/instanceManagement"

	"Engee-Game/utils"
)

func CORSMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusOK)
			c.Writer.Write([]byte(""))
		}

		c.Next()
	}
}

func Serve(port string) {
	router := gin.Default()

	router.Use(CORSMiddleWare())

	router.POST("/games", postGame)
	router.PUT("/games/:id/start", startGame)
	router.PUT("/games/:id/pause", pauseGame)
	router.PUT("/games/:id/reset", resetGame)
	router.PUT("/games/:id/rules", updateGameRules)
	router.DELETE("/games/:id/players/:id", removePlayer)
	router.DELETE("/games/:id", deleteGame)

	router.Run(":" + port)
}

func postGame(c *gin.Context) {
	w := c.Writer

	reqBody := utils.ProcessRequestMessage(c)
	if reqBody == nil {
		return
	}

	err := im.CreateNewInstance(string(reqBody))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Creating game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after creating game: %v", err)
		return
	}
}

func startGame(c *gin.Context) {
	w := c.Writer
	ids := utils.GetRequestIDs(c.Request)

	err := im.StartInstance(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Starting game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after starting game: %v", err)
		return
	}
}

func pauseGame(c *gin.Context) {
	w := c.Writer
	ids := utils.GetRequestIDs(c.Request)

	err := im.PauseInstance(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to pause/unpause game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Pausing/Unpausing game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after pausing game: %v", err)
		return
	}
}

func resetGame(c *gin.Context) {
	w := c.Writer
	ids := utils.GetRequestIDs(c.Request)

	err := im.ResetInstance(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reset game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Resetting game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after resetting game: %v", err)
		return
	}
}

func updateGameRules(c *gin.Context) {
	w := c.Writer
	ids := getRequestIDs(c)

	reqBody := utils.ProcessRequestMessage(c)
	if reqBody == nil {
		return
	}

	err := im.SetInstanceRules(ids[0], string(reqBody))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update game rules: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating game rules: %v", err)
		return
	}
}

func removePlayer(c *gin.Context) {
	w := c.Writer
	ids := utils.GetRequestIDs(c.Request)

	err := im.RemovePlayerFromInstance(ids[0], ids[1])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove player from game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Removing player from game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after removing player: %v", err)
		return
	}
}

func deleteGame(c *gin.Context) {
	w := c.Writer
	ids := getRequestIDs(c)

	err := im.DeleteInstance(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Deleting game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after deleting game: %v", err)
		return
	}
}

func sendReply(w http.ResponseWriter, msg string, code int) error {
	w.WriteHeader(code)
	_, err := w.Write([]byte(msg))
	if err != nil {
		http.Error(w, "Could not write response", http.StatusInternalServerError)
		return fmt.Errorf("could not write response: %v", err)
	}

	return nil
}

func getRequestIDs(c *gin.Context) []string {
	request := c.Request

	splitPath := strings.Split(request.URL.Path, "/")

	ids := make([]string, 0)

	for _, pathPart := range splitPath {
		_, err := uuid.Parse(pathPart)
		if err == nil {
			ids = append(ids, pathPart)
		}
	}

	return ids
}
