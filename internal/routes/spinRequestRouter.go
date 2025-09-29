package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/substrate-cli/api-server/internal/handlers"
	"github.com/substrate-cli/api-server/internal/middlewares"
)

func registerSpinRequest(router *gin.RouterGroup) {
	router.POST("/spin-request", middlewares.RequestLogger(), middlewares.ValidateUser, handlers.InitiateRequest) // currentuser, prompt
}
