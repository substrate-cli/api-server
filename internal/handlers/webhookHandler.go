package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/substrate-cli/api-server/internal/helpers"
	"github.com/substrate-cli/api-server/internal/utils"
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
	magenta := color.New(color.FgHiMagenta).SprintFunc()
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

	var data map[string]any
	err = json.Unmarshal([]byte(payload.Stream), &data)
	if err != nil {
		log.Println("Unable to process stream")
	}
	url := fmt.Sprintf("http://localhost:%d", int(data["appPort"].(float64)))
	stm := fmt.Sprintf("app is running on %s", magenta(url))
	log.Println(stm)
	// log.Println("app is running on port http://localhost:", int(data["appPort"].(float64)))
	utils.StopLoader()
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

	green := color.New(color.FgGreen).SprintFunc()
	log.Println("broadcasting message...")
	log.Println("type: ", payload.Type)
	log.Println("Stream Received: ", green(payload.Stream))

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

func ErrorAction(context *gin.Context) {
	type Payload struct {
		ARCHIVE_KEY string `json:"ARCHIVE_KEY`
		Type        string `json:"type"`
		Stream      string `json:"stream"`
		Status      string `json:"status"`
		Error       string `json:"error"`
	}

	var payload Payload
	err := context.ShouldBindJSON(&payload)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	red := color.New(color.FgRed).SprintFunc()
	log.Println("broadcasting message...")
	log.Println("type: ", payload.Type)
	log.Println("Stream Received: ", red(payload.Stream))

	message, err := json.Marshal(gin.H{
		"ARCHIVE_KEY":   "streamchat",
		"broadCastType": payload.Type,
		"stream":        payload.Stream,
		"status":        payload.Status,
		"error":         payload.Error,
	})

	if err != nil {
		log.Println("Problem parsing json in request")
		context.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}
	log.Println(payload.Error)

	log.Println("Proceeding to broadcast precheck status")
	broadcastMessage(message)
	utils.StopLoader()
	log.Println("restarting server...")
	context.JSON(http.StatusOK, gin.H{"message": "broadcasting precheck completed."})
	helpers.Selector() //restarting selector
}
