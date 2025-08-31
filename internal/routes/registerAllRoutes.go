package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sshfz/api-server-substrate/internal/handlers"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")

	router.GET("/ws", func(ctx *gin.Context) {
		handlers.HandleWS(ctx.Writer, ctx.Request)
	})
	codeGenerationWebhook(api)
	registerSpinRequest(api)
}
