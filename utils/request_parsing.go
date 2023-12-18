package utils

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetRequestIDs(request *http.Request) []string {
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

func ProcessRequestMessage(c *gin.Context) []byte {
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		http.Error(c.Writer, "Failed to read request body", http.StatusBadRequest)
		log.Printf("[Error] Reading request body: %v", err)
		return nil
	}

	return reqBody
}
