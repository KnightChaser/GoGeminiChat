package modelUsages

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"gogeminichat/apikey"
	"gogeminichat/module/responseProcess"
	"gogeminichat/module/utility"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func ChatSessionStreaming(genAIModelName string) {

	apiKeyString := apikey.GetGoogleGenAIAPIKey()

	// Colorful text print on the console
	context := context.Background()
	yellowColorBoldPrint := color.New(color.FgYellow, color.Bold)
	cyanColorBoldPrint := color.New(color.FgCyan, color.Bold)
	whiteColorItalicPrint := color.New(color.FgWhite, color.Italic)
	whiteColorBoldPrint := color.New(color.FgWhite, color.Bold)

	// Ready to bring the model
	fmt.Printf("Obtained an API KEY: %s\n", whiteColorBoldPrint.Sprint(apiKeyString))
	client, err := genai.NewClient(context, option.WithAPIKey(apiKeyString))
	if err != nil {
		log.Panic(err)
	}
	defer client.Close()

	// An interaction count
	chatSessionInteractionCount := 0

	// Prepare the model with disabling safety regulation that may harm the user experience
	model := client.GenerativeModel(genAIModelName)
	chatSession := model.StartChat()
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
	}
	chatSession.History = []*genai.Content{}

	fmt.Printf("Chat session with Gemini AI Model(%s) started.\n", genAIModelName)

	// Start chat session endlessly. User -> Model -> User -> Model
	for {
		// Receive user
		yellowColorBoldPrint.Println(" - user")
		fmt.Print("> ")
		var question string
		reader := bufio.NewReader(os.Stdin)
		question, error := reader.ReadString('\n')
		if error != nil {
			log.Fatal(error)
		}

		// A trigger to exit or save
		if strings.TrimRight(question, "\r\n") == "exit" {
			fmt.Println("(Exiting) Goodbye.")
			break
		} else if strings.TrimRight(question, "\r\n") == "save" {
			fmt.Println("The chat history will be saved into a file.")
			fmt.Print("filename> ")
			var filename string
			reader := bufio.NewReader(os.Stdin)
			filename, error := reader.ReadString('\n')
			if error != nil {
				log.Fatal(error)
			}
			saveChatSessionHistoryIntoJSON(chatSession, filename)
			break
		}

		chatSessionInteractionCount += 1
		// Counting question token
		tokenQtyResponse, error := model.CountTokens(context, genai.Text(question))
		if error != nil {
			log.Fatal(err)
		}
		whiteColorItalicPrint.Printf("...Prompt length: %d tokens\n", tokenQtyResponse.TotalTokens)
		whiteColorItalicPrint.Printf("...User-Model interaction: %d times\n", chatSessionInteractionCount)

		responseIterative := chatSession.SendMessageStream(context, genai.Text(question))
		var answer string

		// Print response from the model via streaming
		cyanColorBoldPrint.Printf(" - %s\n", genAIModelName)
		for {
			response, error := responseIterative.Next()
			if error == iterator.Done {
				break
			} else if error != nil {
				log.Fatal(error)
			}

			result, _ := responseProcess.GetGeminiAITextOnlyResponseStruct(question, response)
			answerIterationFragment := result.Response[0]
			fmt.Print(responseProcess.BoldifyTextInMarkdownRule(answerIterationFragment))
			answer = answer + answerIterationFragment
		}

		fmt.Println()
		responseProcess.AddMessageToChatSessionHistory(chatSession, "user", question)
		responseProcess.AddMessageToChatSessionHistory(chatSession, "model", answer)
	}

}

type ChatHistoryUnit struct {
	Message string
	Role    string
}

type ChatHistory struct {
	Content []ChatHistoryUnit
}

func saveChatSessionHistoryIntoJSON(chatSession *genai.ChatSession, filename string) {
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

	fmt.Printf("Chat history saved to %s as JSON.\n", filepath)
}
