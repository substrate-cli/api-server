package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/manifoldco/promptui"
	"github.com/sshfz/api-server-substrate/cmd/app/connections"
	"github.com/sshfz/api-server-substrate/internal/routes"
)

func main() {
	router := gin.Default()

	connections.InitRabbitMQ()
	connections.InitRedis()

	// Defining the root route
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://your-frontend.com"},
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
		Addr:    ":8080",
		Handler: router,
	}

	stopChan := make(chan struct{})
	startHeartbeat(stopChan)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		log.Println("Server running on http://localhost:8080")
	}()

	selector()

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

func selector() {
	prompt := promptui.Select{
		Label: "Choose your model...",
		Items: []string{"Anthropic Claude 4.1", "GPT-5", "Gemini Flash 2.0"},
	}

	// Run the picker
	index, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	fmt.Printf("#%d: %s\n", index, result)

	descPrompt := promptui.Prompt{
		Label: "Describe here...",
	}

	description, err := descPrompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	stop := make(chan bool)
	go showLoader(stop)

	// Simulate some work (e.g., API call)
	time.Sleep(3 * time.Second)

	// Stop loader
	stop <- true

	log.Println(description)
}

func showLoader(stop chan bool) {
	chars := `|/-\`
	i := 0
	for {
		select {
		case <-stop:
			return
		default:
			fmt.Printf("\rThinking %c", chars[i%len(chars)])
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}
