package modelUsages

import (
	"bufio"
	"context"
	"fmt"
	"gogeminichat/apikey"
	"gogeminichat/module/responseProcess"
	"gogeminichat/module/utility"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func ChatSessionStreaming(genAIModelName string) {

	apiKeyString := apikey.GetGoogleGenAIAPIKey()

	// Ready to bring the model
	context := context.Background()
	fmt.Printf("Obtained an API KEY: %s\n", utility.WhiteColorBoldPrint.Sprint(apiKeyString))
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
		utility.YellowColorBoldPrint.Println(" - user")
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
			utility.CyanColorBoldPrint.Print("filename to save...> ")
			var filename string
			reader := bufio.NewReader(os.Stdin)
			filename, error := reader.ReadString('\n')
			if error != nil {
				log.Fatal(error)
			}
			responseProcess.SaveChatSessionHistoryIntoJSON(chatSession, filename)
			break
		}

		chatSessionInteractionCount += 1
		// Counting question token
		tokenQtyResponse, error := model.CountTokens(context, genai.Text(question))
		if error != nil {
			log.Fatal(err)
		}
		utility.WhiteColorItalicPrint.Printf("...Prompt length: %d tokens\n", tokenQtyResponse.TotalTokens)
		utility.WhiteColorItalicPrint.Printf("...User-Model interaction: %d times\n", chatSessionInteractionCount)

		responseIterative := chatSession.SendMessageStream(context, genai.Text(question))
		var answer string

		// Print response from the model via streaming
		utility.CyanColorBoldPrint.Printf(" - %s\n", genAIModelName)
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
