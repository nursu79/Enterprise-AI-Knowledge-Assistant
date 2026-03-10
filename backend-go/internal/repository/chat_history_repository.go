package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nursu79/go-production-api/internal/repository/storage"
)

type ChatHistoryRepository interface {
	CreateChatHistory(ctx context.Context, arg storage.CreateChatHistoryParams) (storage.ChatHistory, error)
	GetChatHistoryByUserID(ctx context.Context, arg storage.GetChatHistoryByUserIDParams) ([]storage.ChatHistory, error)
}

type chatHistoryRepository struct {
	db      *pgxpool.Pool
	queries *storage.Queries
}

func NewChatHistoryRepository(db *pgxpool.Pool) ChatHistoryRepository {
	return &chatHistoryRepository{
		db:      db,
		queries: storage.New(db),
	}
}

func (r *chatHistoryRepository) CreateChatHistory(ctx context.Context, arg storage.CreateChatHistoryParams) (storage.ChatHistory, error) {
	return r.queries.CreateChatHistory(ctx, arg)
}

func (r *chatHistoryRepository) GetChatHistoryByUserID(ctx context.Context, arg storage.GetChatHistoryByUserIDParams) ([]storage.ChatHistory, error) {
	return r.queries.GetChatHistoryByUserID(ctx, arg)
}
