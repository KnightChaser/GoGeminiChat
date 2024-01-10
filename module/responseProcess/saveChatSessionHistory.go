package responseProcess

import (
	"encoding/json"
	"fmt"
	"gogeminichat/module/utility"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

type ChatHistoryUnit struct {
	Message string
	Role    string
}

type ChatHistory struct {
	Content []ChatHistoryUnit
}

func SaveChatSessionHistoryIntoJSON(chatSession *genai.ChatSession, filename string) {
	// Initialize an empty Content slice
	var chatSessionContent ChatHistory

	// Calculate the halfway point of the iteration
	// I don't know why... but *genai.ChatSession contains saves the same history twice.
	halfway := len(chatSession.History) / 2

	for messageChunkIndex := 0; messageChunkIndex < halfway; messageChunkIndex++ {
		// Assuming Parts is a slice of strings, join them into a single string
		message := fmt.Sprintf("%v", chatSession.History[messageChunkIndex].Parts[0])
		role := chatSession.History[messageChunkIndex].Role

		// Append new ChatHistoryUnit to Content
		chatSessionContent.Content = append(chatSessionContent.Content, ChatHistoryUnit{
			Message: message,
			Role:    role,
		})
	}

	JSONconverted, err := json.MarshalIndent(chatSessionContent, "", " ")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Ready to save the log into the file
	filename = strings.TrimRight(filename, "\r\n") // Ignoring inserted newline characters (\r\n) will occur an error
	currentWorkingDirectory := utility.GetCurrentExecutablePath()
	filepath := fmt.Sprintf("%s/userChatHistory/%s.json", currentWorkingDirectory, filename)

	// Open file for writing
	JSONExportedFile, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer JSONExportedFile.Close()

	// Write JSON data to file
	_, err = JSONExportedFile.Write(JSONconverted)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Chat history saved to %s as JSON.\n", utility.CyanColorBoldPrint.Sprint(filepath))
}
