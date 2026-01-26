package utils

import (
	"context"
	"errors"
	"os"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
}

// Menginisialisasi client Gemini menggunakan GEMINI_API_KEY.
func NewGeminiClient() *GeminiClient {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		panic("GEMINI_API_KEY is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		panic(err)
	}

	return &GeminiClient{client: client}
}

// Generate menghasilkan teks dari model Gemini berdasarkan prompt.
func (g *GeminiClient) Generate(prompt string) (string, error) {
	ctx := context.Background()

	resp, err := g.client.Models.GenerateContent(ctx, "gemini-3-flash-preview", genai.Text(prompt), nil)
	if err != nil {
		return "", err
	}

	if resp.Text() == "" {
		return "", errors.New("no output from Gemini model")
	}

	return resp.Text(), nil
}
