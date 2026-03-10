package usecase

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nursu79/go-production-api/internal/repository"
	"github.com/nursu79/go-production-api/internal/repository/storage"
)

type ChatHistoryUsecase interface {
	SaveChatHistory(ctx context.Context, userID pgtype.UUID, query, aiResponse string, retrievedContext []byte) (storage.ChatHistory, error)
	GetChatHistory(ctx context.Context, userID pgtype.UUID, limit, offset int32) ([]storage.ChatHistory, error)
}

type chatHistoryUsecase struct {
	repo repository.ChatHistoryRepository
}

func NewChatHistoryUsecase(repo repository.ChatHistoryRepository) ChatHistoryUsecase {
	return &chatHistoryUsecase{
		repo: repo,
	}
}

func (u *chatHistoryUsecase) SaveChatHistory(ctx context.Context, userID pgtype.UUID, query, aiResponse string, retrievedContext []byte) (storage.ChatHistory, error) {
	arg := storage.CreateChatHistoryParams{
		UserID:           userID,
		Query:            query,
		AiResponse:       aiResponse,
		RetrievedContext: retrievedContext,
	}
	return u.repo.CreateChatHistory(ctx, arg)
}

func (u *chatHistoryUsecase) GetChatHistory(ctx context.Context, userID pgtype.UUID, limit, offset int32) ([]storage.ChatHistory, error) {
	arg := storage.GetChatHistoryByUserIDParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}
	return u.repo.GetChatHistoryByUserID(ctx, arg)
}
