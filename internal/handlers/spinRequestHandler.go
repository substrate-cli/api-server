package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sshfz/api-server-substrate/cmd/app/connections"
	"github.com/sshfz/api-server-substrate/internal/db"
)

func InitiateRequest(context *gin.Context) {
	type SpinRequest struct {
		UserId string `json:"userid"`
		Prompt string `json:"prompt"`
	}

	var spinRequest SpinRequest

	err := context.ShouldBindJSON(&spinRequest)
	if err != nil {
		log.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	routingKey := "spin.create"

	type amqpReq struct {
		UserId  string
		Message string
		Prompt  string
	}

	var req amqpReq = amqpReq{
		UserId:  spinRequest.UserId,
		Message: "spin-project",
		Prompt:  spinRequest.Prompt,
	}

	err = connections.PublishSpinRequest(req, routingKey)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "published spin, spin init requested."})
}

func UpdateRequest(context *gin.Context) {
	type SpinRequest struct {
		UserId string `json: userid`
		Prompt string `json: prompt`
	}

	cluster := context.Param("cluster")
	if cluster == "" {
		context.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	clusterExists, err := db.ReadValueFromKey(cluster)
	if err != nil {
		log.Println("error reading cluster")
		context.JSON(http.StatusForbidden, gin.H{"error": "Cluster not found"})
		return
	}
	if clusterExists == "" {
		context.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}
	if clusterExists != "running" {
		context.JSON(http.StatusForbidden, gin.H{"error": "Cluster not found"})
		return
	}

	log.Println("cluster found =>", cluster)

	var spinRequest SpinRequest

	err = context.ShouldBindJSON(&spinRequest)
	if err != nil {
		log.Println("Bad Request", err)
		context.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	type amqpReq struct {
		UserId  string
		Message string
		Prompt  string
		Cluster string
	}

	routingKey := "spin.update"

	var req amqpReq = amqpReq{
		UserId:  spinRequest.UserId,
		Message: "spin-update",
		Prompt:  spinRequest.Prompt,
		Cluster: cluster,
	}

	err = connections.PublishSpinRequest(req, routingKey)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "published update, spin update requested."})
}
