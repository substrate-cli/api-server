package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CodeGenerationComplete(context *gin.Context) { //this request will only be called from consumer-service --
	type Payload struct {
		ARCHIVE_KEY string `json:"ARCHIVE_KEY`
		Type        string `json:"type"`
		Stream      string `json:"stream"`
		Status      string `json:"status"`
	}

	//save to db
	var payload Payload
	err := context.ShouldBindJSON(&payload)
	if err != nil {
		log.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Notifying UI via WebSocket
	log.Println("broadcasting message...")
	log.Println(payload)
	message, _ := json.Marshal(gin.H{
		"ARCHIVE_KEY":   "streamchat",
		"broadcastType": payload.Type,
		"stream":        payload.Stream,
		"status":        "finished",
	})
	broadcastMessage(message)

	context.JSON(http.StatusOK, gin.H{"message": "broadcasting code generation completed"})
}

func PrecheckAction(context *gin.Context) {
	type Payload struct {
		ARCHIVE_KEY string `json:"ARCHIVE_KEY`
		Type        string `json:"type"`
		Stream      string `json:"stream"`
		Status      string `json:"status"`
	}

	var payload Payload
	err := context.ShouldBindJSON(&payload)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	log.Println("broadcasting message...")
	log.Println(payload)

	message, err := json.Marshal(gin.H{
		"ARCHIVE_KEY":   "streamchat",
		"broadCastType": payload.Type,
		"stream":        payload.Stream,
		"status":        payload.Status,
	})

	if err != nil {
		log.Println("Problem parsing json in request")
		context.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}
	log.Println("Proceeding to broadcast precheck status")
	broadcastMessage(message)
	context.JSON(http.StatusOK, gin.H{"message": "broadcasting precheck completed."})
}
