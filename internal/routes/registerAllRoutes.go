package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/substrate-cli/api-server/internal/handlers"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")

	router.GET("/ws", func(ctx *gin.Context) {
		handlers.HandleWS(ctx.Writer, ctx.Request)
	})
	codeGenerationWebhook(api)
	registerSpinRequest(api)
}
