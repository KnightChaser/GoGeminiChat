package apikey

import (
	"fmt"
	"log"
	"os"
)

// Grab the API key for the service usage, reading API key from the designated file.
func GetGoogleGenAIAPIKey() string {
	// Setup API key (the given API file will contain text of API key)

	geminiAPIkeyFilePath := "apikey/apikey.txt"

	geminiAPIKey, err := os.ReadFile(geminiAPIkeyFilePath)
	if err != nil {
		fmt.Printf("Failed to load the API key (Expected API key: %v)\n", geminiAPIkeyFilePath)
		log.Panic(err)
		return "ERR"
	}

	// Convert geminiAPIKey from []byte to string
	apiKeyString := string(geminiAPIKey)

	return apiKeyString
}
