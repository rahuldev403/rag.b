package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rahuldev403/forgeflow/internal/config"
)

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func AnalyzePayload(prompt string, payloadData string) (string, error) {
	apiKey := config.GetEnv("GEMINI_API_KEY", "")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not set")
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-3.7-flash:generateContent?key=" + apiKey

	fullPrompt := fmt.Sprintf("%s\n\nData:\n%s", prompt, payloadData)

	reqBody := GeminiRequest{
		Contents: []Content{
			{Parts: []Part{{Text: fullPrompt}}},
		},
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("network error calling AI API: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var aiResponse GeminiResponse
	if err := json.Unmarshal(bodyBytes, &aiResponse); err != nil {
		return "", fmt.Errorf("failed to parse AI response:%v", err)
	}

	if len(aiResponse.Candidates) > 0 && len(aiResponse.Candidates[0].Content.Parts) > 0 {
		return aiResponse.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("AI returned an empty response")
}
