package responseProcess

import (
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

type GeminiTextResponseSafetyRating struct {
	Probability           string
	BlockedBySafetyPolicy bool
}

type GeminiTextResponseStructure struct {
	Question      string
	ResponseCount uint32
	Role          string
	Response      []string
	FininshReason string
	TokenCount    uint32
	SafetyRating  map[string]GeminiTextResponseSafetyRating
}

// Using the custom response object of "GeminiTextResponseStructure", preprocess the raw data from the API to the normal Go data structure with Go datatypes.
func GetGeminiAITextOnlyResponseStruct(question string, response *genai.GenerateContentResponse) (GeminiTextResponseStructure, error) {
	var responseStructure GeminiTextResponseStructure
	responseStructure.Question = question
	for _, candidates := range response.Candidates {
		if candidates.Content != nil {
			// response text
			for _, part := range candidates.Content.Parts {
				responseStructure.ResponseCount += 1
				responseStructure.Response = append(responseStructure.Response, fmt.Sprintf("%s", part))
			}

			// metadata
			responseStructure.Role = candidates.Content.Role
			responseStructure.FininshReason = fmt.Sprintf("%s", candidates.FinishReason)
			responseStructure.TokenCount = uint32(candidates.TokenCount)
			responseStructure.SafetyRating = make(map[string]GeminiTextResponseSafetyRating)
			for _, safetyData := range candidates.SafetyRatings {
				safetyCategory := fmt.Sprintf("%s", safetyData.Category)

				// Check if the map entry exists, create it if not
				if _, ok := responseStructure.SafetyRating[safetyCategory]; !ok {
					responseStructure.SafetyRating[safetyCategory] = GeminiTextResponseSafetyRating{}
				}

				// Create an instance of GeminiTextResponseSafetyRating
				safetyRating := responseStructure.SafetyRating[safetyCategory]

				// Assign values to the struct fields
				safetyRating.BlockedBySafetyPolicy = safetyData.Blocked
				safetyRating.Probability = fmt.Sprintf("%s", safetyData.Probability)

				// Update the map entry
				responseStructure.SafetyRating[safetyCategory] = safetyRating
			}

		} else {
			return GeminiTextResponseStructure{}, fmt.Errorf("Failed to receive data from Google Gen. API")
		}
	}
	return responseStructure, nil
}