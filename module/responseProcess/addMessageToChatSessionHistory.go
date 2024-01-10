package responseProcess

import "github.com/google/generative-ai-go/genai"

func AddMessageToChatSessionHistory(chatSession *genai.ChatSession, role string, text string) {
	message := &genai.Content{
		Parts: []genai.Part{
			genai.Text(text),
		},
		Role: role,
	}
	chatSession.History = append(chatSession.History, message)
}
