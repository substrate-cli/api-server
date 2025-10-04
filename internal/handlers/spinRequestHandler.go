package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/substrate-cli/api-server/cmd/app/connections"
	"github.com/substrate-cli/api-server/internal/db"
	"github.com/substrate-cli/api-server/internal/helpers"
	"github.com/substrate-cli/api-server/internal/utils"
)

func InitiateRequest(context *gin.Context) {
	val, err := db.ReadValueFromKey(utils.GetDefaultUser())
	if val == "processing" {
		context.JSON(http.StatusForbidden, gin.H{"err": "Substrate is busy baking one of your requests. Please check back in a little while."})
		return
	}
	db.SaveRedis(utils.GetDefaultUser(), "processing")
	type SpinRequest struct {
		UserId      string `json:"userid"`
		Prompt      string `json:"prompt"`
		ClusterName string `json:"clustername"`
		Model       string `json:"model"`
	}

	var spinRequest SpinRequest

	err = context.ShouldBindJSON(&spinRequest)
	if err != nil {
		db.SaveRedis(utils.GetDefaultUser(), "failed")
		log.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	routingKey := "spin.create"

	if strings.ReplaceAll(spinRequest.Prompt, " ", "") == "" {
		db.SaveRedis(utils.GetDefaultUser(), "failed")
		context.JSON(http.StatusBadRequest, gin.H{"err": "prompt should not be empty"})
		return
	}

	type amqpReq struct {
		UserId      string
		Message     string
		Prompt      string
		ClusterName string
		Model       string
	}

	clusterName := strings.TrimSpace(spinRequest.ClusterName)
	clusterName = strings.ReplaceAll(clusterName, " ", "-")

	var req amqpReq = amqpReq{
		UserId:      utils.GetDefaultUser(),
		Message:     "spin-project",
		Prompt:      spinRequest.Prompt,
		ClusterName: clusterName,
		Model:       spinRequest.Model,
	}

	if len(clusterName) != 0 && helpers.CheckIfDirExists(req.ClusterName) {
		db.SaveRedis(utils.GetDefaultUser(), "failed")
		log.Println(clusterName, "Directory already exists, try a different name")
		context.JSON(http.StatusForbidden, gin.H{"error": "Cluster with this name already exists, try a different name"})
		return
	}

	err = connections.PublishSpinRequest(req, routingKey)
	if err != nil {
		db.SaveRedis(utils.GetDefaultUser(), "failed")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "published spin, spin init requested."})
}
