package helpers

import (
	"log"
	"os"
	"path/filepath"
)

func CheckIfDirExists(clusterName string) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("unable to retrieve home directory")
		return false
	}
	rootProjectPath := filepath.Join(homeDir, "Desktop", "substrate-home", clusterName)

	info, err := os.Stat(rootProjectPath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return info.IsDir()
}
