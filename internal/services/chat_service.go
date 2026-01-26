package services

import (
	"fmt"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/utils"
)

type ChatService struct {
	repo   repository.ChatRepository
	gemini *utils.GeminiClient
}

func NewChatService(r repository.ChatRepository, g *utils.GeminiClient) *ChatService {
	return &ChatService{
		repo:   r,
		gemini: g,
	}
}

func (s *ChatService) Chat(req dto.ChatRequest) (string, error) {
	history, _ := s.repo.GetHistory(req.SessionID)

	prompt := "Kamu adalah chatbot CS.\n\n"
	for _, h := range history {
		prompt += "User: " + h.User + "\n"
		prompt += "Bot: " + h.Bot + "\n"
	}
	prompt += "User: " + req.Message + "\n"

	answer, err := s.gemini.Generate(prompt)
	if err != nil {
		fmt.Println("Gemini Generate Error:", err)
		return "", err
	}

	_ = s.repo.Save(req.SessionID, req.Message, answer)
	return answer, nil
}

func (s *ChatService) Repo() repository.ChatRepository {
	return s.repo
}

// Expose Gemini client
func (s *ChatService) GetGemini() *utils.GeminiClient {
	return s.gemini
}
