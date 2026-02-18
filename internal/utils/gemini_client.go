package utils

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
}

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

func (g *GeminiClient) Generate(prompt string) (string, error) {
	models := []string{
		"gemma-3-12b-it",
		"gemma-2-27b-it",
		"gemini-1.5-flash",
	}

	// Try each model
	for modelIdx, modelName := range models {
		if modelIdx > 0 {
			log.Printf("Fallback ke model: %s", modelName)
		}

		maxRetries := 5
		baseDelay := 2 * time.Second

		for attempt := 0; attempt < maxRetries; attempt++ {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			resp, err := g.client.Models.GenerateContent(
				ctx,
				modelName,
				genai.Text(prompt),
				nil,
			)
			cancel()

			if err != nil {
				errMsg := err.Error()

				// Check error type
				is503 := strings.Contains(errMsg, "503") ||
					strings.Contains(errMsg, "UNAVAILABLE") ||
					strings.Contains(errMsg, "overloaded")

				isRateLimit := strings.Contains(errMsg, "RESOURCE_EXHAUSTED") ||
					strings.Contains(errMsg, "quota exceeded") ||
					strings.Contains(errMsg, "429")

				isModelNotFound := strings.Contains(errMsg, "NOT_FOUND")

				// if model not found, langsung coba model berikutnya
				if isModelNotFound {
					log.Printf("Model %s tidak tersedia", modelName)
					break
				}

				shouldRetry := is503 || isRateLimit

				if shouldRetry && attempt < maxRetries-1 {
					// Exponential backoff: 2s, 4s, 8s, 16s, 32s
					delay := baseDelay * time.Duration(1<<uint(attempt))
					// Jitter untuk avoid thundering herd
					jitter := time.Duration(time.Now().UnixNano()%1000) * time.Millisecond
					totalDelay := delay + jitter

					log.Printf("Retry %d/%d dalam %v... (error: %s)",
						attempt+1, maxRetries, totalDelay, errMsg)

					time.Sleep(totalDelay)
					continue
				}

				// If max retry, coba model berikutnya
				if attempt == maxRetries-1 && modelIdx < len(models)-1 {
					log.Printf("Max retries untuk %s, coba model lain...", modelName)
					break
				}

				// If model terakhir dan sudah max retry
				if attempt == maxRetries-1 && modelIdx == len(models)-1 {
					return "", errors.New("semua model gagal setelah beberapa kali retry")
				}

				continue
			}

			// Success
			if resp.Text() == "" {
				return "", errors.New("no output from Gemini model")
			}

			if attempt > 0 || modelIdx > 0 {
				log.Printf("Berhasil dengan %s setelah %d retry", modelName, attempt)
			}

			return resp.Text(), nil
		}
	}

	return "", errors.New("semua model gagal")
}
