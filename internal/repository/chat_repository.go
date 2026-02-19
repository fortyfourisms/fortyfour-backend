package repository

import (
	"fmt"
	"fortyfour-backend/internal/dto"
)

type ChatRepository interface {
	GetHistory(sessionID string) ([]dto.ChatHistory, error)
	Save(sessionID, userMsg, botMsg string) error
}

type InMemoryChatRepo struct {
	data map[string][]dto.ChatHistory
}

func NewInMemoryChatRepo() *InMemoryChatRepo {
	return &InMemoryChatRepo{
		data: make(map[string][]dto.ChatHistory),
	}
}

func (r *InMemoryChatRepo) GetHistory(sessionID string) ([]dto.ChatHistory, error) {
	return r.data[sessionID], nil
}

func (r *InMemoryChatRepo) Save(sessionID, userMsg, botMsg string) error {
	r.data[sessionID] = append(r.data[sessionID], dto.ChatHistory{
		User: userMsg,
		Bot:  botMsg,
	})
	return nil
}

func (r *InMemoryChatRepo) DeleteSession(sessionID string) error {
	if _, exists := r.data[sessionID]; !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}
	delete(r.data, sessionID)
	return nil
}
