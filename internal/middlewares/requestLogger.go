package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()
		duration := time.Since(start)
		log.Printf("[%s] %s %s %d in %v",
			context.Request.Method,
			context.Request.Host,
			context.Request.URL.Path,
			context.Writer.Status(),
			duration,
		)
		context.Next()
	}
}
