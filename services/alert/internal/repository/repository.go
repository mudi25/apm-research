package repository

import (
	"context"
	"research-apm/services/alert/internal/entity"
	"research-apm/services/alert/internal/repository/internal/repository"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/go-telegram/bot"
)

type Repository interface {
	// get apm alert
	GetAlert(ctx context.Context) ([]entity.Alert, error)

	// get apm alert
	SendTelegram(ctx context.Context, data entity.Alert) error
}

func NewRepository(esClient *elasticsearch.Client, telegramBot *bot.Bot, chatID string) Repository {
	return repository.NewRepository(esClient, telegramBot, chatID)
}
