package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sshfz/api-server-substrate/internal/handlers"
	"github.com/sshfz/api-server-substrate/internal/middlewares"
)

func registerSpinRequest(router *gin.RouterGroup) {
	router.POST("/spin-request", middlewares.RequestLogger(), middlewares.ValidateUser, handlers.InitiateRequest) // currentuser, prompt
}
