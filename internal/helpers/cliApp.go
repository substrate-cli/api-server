package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
	"github.com/substrate-cli/api-server/cmd/app/connections"
	"github.com/substrate-cli/api-server/internal/utils"
)

func Selector() {
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
			if len(strings.TrimSpace(s)) != 0 && CheckIfDirExists(s) {
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

func getAPIKey() (string, error) {
	// Configure readline with better settings
	config := &readline.Config{
		Prompt:          "Enter API Key: ",
		HistoryFile:     "", // No history file for security
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		// EnableMask:      true,
		// MaskRune:        '*',
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
