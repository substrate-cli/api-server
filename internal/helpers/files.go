package helpers

import (
	"log"
	"os"
	"path/filepath"

	"github.com/substrate-cli/api-server/internal/utils"
)

func CheckIfDirExists(clusterName string) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("unable to retrieve home directory")
		return false
	}
	rootProjectPath := ""
	// if utils.GetBundle() == "docker" {
	// 	return false
	// }
	if utils.GetBundle() == "docker" {
		homeDir = "/apps"
		rootProjectPath = filepath.Join(homeDir, "substrate-home", clusterName)
	} else {
		rootProjectPath = filepath.Join(homeDir, "Desktop", "substrate-home", clusterName)
	}

	info, err := os.Stat(rootProjectPath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return info.IsDir()
}
