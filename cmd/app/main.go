package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/substrate-cli/api-server/cmd/app/connections"
	"github.com/substrate-cli/api-server/internal/helpers"
	"github.com/substrate-cli/api-server/internal/routes"
	"github.com/substrate-cli/api-server/internal/utils"
)

func main() {
	router := gin.Default()

	connections.InitRabbitMQ()
	connections.InitRedis()

	origins := utils.GetSafeOrigins()
	parts := strings.Split(origins, ",")
	safeOrigins := make([]string, 0, len(parts))
	for _, o := range parts {
		safeOrigins = append(safeOrigins, strings.TrimSpace(o))
	}

	// Defining the root route
	router.Use(cors.New(cors.Config{
		AllowOrigins:     safeOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "substrate release - 1.0.0")
	})

	routes.RegisterRoutes(router)
	fig := figure.NewFigure("substrate-cli", "shadow", true)
	color.Set(color.FgRed)
	fig.Print()
	color.Unset()

	srv := &http.Server{
		Addr:    ":" + utils.GetPort(),
		Handler: router,
	}

	stopChan := make(chan struct{})
	startHeartbeat(stopChan)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		log.Println("Server running on http://localhost:" + utils.GetPort())
	}()

	mode := utils.GetMode()
	if mode == "cli" {
		helpers.Selector()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Shutting down server...")
	close(stopChan)
	// Graceful shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting gracefully")
}

func clearLineAndLog(message string) {
	// Clear current line, print log message, then move to new line
	fmt.Print("\r\033[K")
	log.Print(message)
}
