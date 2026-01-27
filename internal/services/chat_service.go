package services

import (
	"fmt"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/utils"
	"strings"
)

type ChatService struct {
	repo           repository.ChatRepository
	perusahaanRepo *repository.PerusahaanRepository
	gemini         *utils.GeminiClient
}

func NewChatService(
	r repository.ChatRepository,
	p *repository.PerusahaanRepository,
	g *utils.GeminiClient,
) *ChatService {
	return &ChatService{
		repo:           r,
		perusahaanRepo: p,
		gemini:         g,
	}
}

// ambil DATA dari DB
func (s *ChatService) BuildDataContext(intent string) string {
	switch intent {

	case "latest_perusahaan":
		data, err := s.perusahaanRepo.GetLatest(5)
		if err != nil || len(data) == 0 {
			return "Data perusahaan tidak tersedia di sistem."
		}

		var b strings.Builder
		for _, p := range data {
			b.WriteString(fmt.Sprintf(
				"- %s (sektor: %s, dibuat %s)\n",
				p.NamaPerusahaan,
				p.Sektor,
				p.CreatedAt,
			))

		}
		return b.String()

	default:
		return "Data tidak tersedia di sistem."
	}
}

func (s *ChatService) Repo() repository.ChatRepository {
	return s.repo
}

func (s *ChatService) GetGemini() *utils.GeminiClient {
	return s.gemini
}
