package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sshfz/api-server-substrate/cmd/app/connections"
)

func InitiateRequest(context *gin.Context) {
	type SpinRequest struct {
		UserId      string `json:"userid"`
		Prompt      string `json:"prompt"`
		ClusterName string `json:"clustername"`
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
		UserId      string
		Message     string
		Prompt      string
		ClusterName string
	}

	clusterName := strings.TrimSpace(spinRequest.ClusterName)
	clusterName = strings.ReplaceAll(clusterName, " ", "-")

	var req amqpReq = amqpReq{
		UserId:      spinRequest.UserId,
		Message:     "spin-project",
		Prompt:      spinRequest.Prompt,
		ClusterName: clusterName,
	}

	err = connections.PublishSpinRequest(req, routingKey)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "published spin, spin init requested."})
}
