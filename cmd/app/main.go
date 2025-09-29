package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/manifoldco/promptui"
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
		selector()
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

func selector() {
	supportedModels := strings.Split(utils.GetSupportedModels(), ",")
	prompt := promptui.Select{
		Label: "Choose your model...",
		Items: supportedModels,
	}

	// Run the picker
	index, model, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}
	model = strings.TrimSpace(model)
	fmt.Printf("#%d: %s\n", index, model)

	clusterName := promptui.Prompt{
		Label: "Cluster Name...",
		Validate: func(s string) error {
			if len(strings.TrimSpace(s)) != 0 && helpers.CheckIfDirExists(s) {
				return errors.New("directory already exists, choose a different name")
			}
			return nil
		},
	}
	cluster, err := clusterName.Run()
	cluster = strings.TrimSpace(cluster)
	cluster = strings.ReplaceAll(cluster, " ", "-")
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	apiKey, _ := getAPIKey()
	description, err := getDesc()

	utils.SetAPIKey(apiKey)

	description = strings.TrimSpace(description)

	utils.StartLoader("thinking")
	user := utils.GetDefaultUser()

	err = publishMessageToConsumer(user, description, cluster, model)
	if err != nil {
		log.Println(err)
		utils.StopLoader()
	}

	// clearLineAndLog(description)
}

func clearLineAndLog(message string) {
	// Clear current line, print log message, then move to new line
	fmt.Print("\r\033[K")
	log.Print(message)
}

func publishMessageToConsumer(user string, prompt string, cluster string, model string) error {
	routingKey := "spin.create"

	type amqpReqCLI struct {
		UserId      string
		Message     string
		Prompt      string
		ApiKey      string
		ClusterName string
		Model       string
	}

	var req amqpReqCLI = amqpReqCLI{
		UserId:      user,
		Message:     "spin-project",
		Prompt:      prompt,
		ApiKey:      *utils.GetAPIKey(),
		ClusterName: cluster,
		Model:       model,
	}

	err := connections.PublishSpinRequest(req, routingKey)
	if err != nil {
		log.Println("there was an error while publishing message to consumer")
		log.Println(err)
		return err
	}

	return nil
}

func getAPIKey() (string, error) {
	// Configure readline with better settings
	config := &readline.Config{
		Prompt:          "Enter API Key: ",
		HistoryFile:     "", // No history file for security
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		return "", err
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err != nil {
			return "", err // Handle Ctrl+C or EOF
		}

		// Validate the input
		apiKey := strings.TrimSpace(input)
		if len(apiKey) == 0 {
			fmt.Println("❌ API key cannot be empty")
			rl.SetPrompt("Enter API Key: ")
			continue // Ask again
		}
		return apiKey, nil
	}
}

func cleanInput(input string) string {
	// Regex to match ANSI escape sequences
	ansiEscape := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiEscape.ReplaceAllString(input, "")
}

func getDesc() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter App description: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		// Validate the input
		desc := strings.TrimSpace(input)
		if len(desc) == 0 {
			fmt.Println("❌ API key cannot be empty")
			continue // Ask again
		}

		return desc, nil
	}
}
