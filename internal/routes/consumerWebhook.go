package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/substrate-cli/api-server/internal/handlers"
	"github.com/substrate-cli/api-server/internal/middlewares"
)

func codeGenerationWebhook(router *gin.RouterGroup) {
	router.POST("/webhook/code-generation", middlewares.RequestLogger(), handlers.CodeGenerationComplete)
	router.POST("/webhook/precheck", middlewares.RequestLogger(), handlers.PrecheckAction)
	router.POST("/webhook/error", middlewares.RequestLogger(), handlers.ErrorAction)
}
