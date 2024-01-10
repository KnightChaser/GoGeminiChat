package main

import (
	"fmt"
	"gogeminichat/module/modelUsages"
)

func main() {
	fmt.Println("Welcome to GoGeminiChat!")
	modelUsages.ChatSessionStreaming("gemini-pro")
}
