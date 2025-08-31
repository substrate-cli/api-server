package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sshfz/api-server-substrate/internal/handlers"
	"github.com/sshfz/api-server-substrate/internal/middlewares"
)

func codeGenerationWebhook(router *gin.RouterGroup) {
	router.POST("/webhook/code-generation", middlewares.RequestLogger(), handlers.CodeGenerationComplete)
	router.POST("/webhook/precheck", middlewares.RequestLogger(), handlers.PrecheckAction)
}
