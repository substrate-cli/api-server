package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type configuration struct {
	port            string
	node            string
	mode            string
	defaultUser     string
	supportedModels string
}

var config *configuration
var apiKey *string

func init() {
	_ = godotenv.Load()

	config = &configuration{
		port:            os.Getenv("PORT"),
		node:            os.Getenv("NODE"),
		mode:            os.Getenv("MODE"),
		defaultUser:     os.Getenv("DEFAULT_USER"),
		supportedModels: os.Getenv("SUPPORTED_MODELS"),
	}
}

func GetPort() string {
	return config.port
}

func SetAPIKey(key string) {
	if config.mode == "cli" {
		apiKey = &key
	} else {
		log.Println("mutating api key not allowed")
	}
}

func GetNode() string {
	return config.node
}

func GetMode() string {
	return config.mode
}

func GetAPIKey() *string {
	return apiKey
}

func GetDefaultUser() string {
	return config.defaultUser
}

func GetSupportedModels() string {
	return config.supportedModels
}
