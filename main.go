package main

import (
	"fmt"
	"gogeminichat/apikey"
)

func main() {
	fmt.Printf("Obtained an API key: %v\n", apikey.GetGoogleGenAIAPIKey())
}
