package middlewares

import "github.com/gin-gonic/gin"

func ValidateUser(context *gin.Context) {
	context.Next()
}
