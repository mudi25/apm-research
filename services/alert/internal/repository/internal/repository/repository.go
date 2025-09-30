package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"research-apm/services/alert/internal/entity"
	"research-apm/services/alert/internal/repository/internal/model"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Repository struct {
	esClient    *elasticsearch.Client
	telegramBot *bot.Bot
	chatID      string
}

func NewRepository(esClient *elasticsearch.Client, telegramBot *bot.Bot, chatID string) *Repository {
	return &Repository{esClient: esClient, telegramBot: telegramBot, chatID: chatID}
}

// get apm alert
func (repo *Repository) GetAlert(ctx context.Context) ([]entity.Alert, error) {
	query := `{
		"_source": [
			"service.name",
			"service.environment",
			"kibana.alert.rule.category",
			"kibana.alert.rule.name",
			"kibana.alert.reason",
			"kibana.alert.status",
			"@timestamp"
		],
		"size": 100,
		"query": {
			"range": {
			"@timestamp": {
				"gte": "now-3h"
			}
			}
		}
	}`
	res, err := repo.esClient.Search(
		repo.esClient.Search.WithContext(ctx),
		repo.esClient.Search.WithIndex(".internal.alerts-observability.apm.alerts-*"),
		repo.esClient.Search.WithBody(strings.NewReader(query)),
		repo.esClient.Search.WithTrackTotalHits(false),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error search: %s", res.String())
	}
	var r struct {
		Hits struct {
			Hits []model.AlertAPMHit `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	alerts := make([]entity.Alert, 0, len(r.Hits.Hits))
	for _, h := range r.Hits.Hits {
		alerts = append(alerts, h.ToEntity())
	}
	return alerts, nil

}

// get apm alert
func (repo *Repository) SendTelegram(ctx context.Context, data entity.Alert) error {
	query := fmt.Sprintf(`{"query": {"term": {"trx_id": "%s"}}}`, data.TrxID)

	countRes, err := repo.esClient.Count(
		repo.esClient.Count.WithContext(ctx),
		repo.esClient.Count.WithIndex("alert-notify-*"),
		repo.esClient.Count.WithBody(strings.NewReader(query)),
		repo.esClient.Count.WithIgnoreUnavailable(true),
	)
	if err != nil {
		return err
	}
	defer countRes.Body.Close()

	if countRes.IsError() {
		return fmt.Errorf("error checking document: %s", countRes.String())
	}

	var c struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(countRes.Body).Decode(&c); err != nil {
		return err
	}

	if c.Count > 0 {
		return nil
	}

	if _, err := repo.telegramBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    repo.chatID,
		Text:      model.NewMessage(data),
		ParseMode: models.ParseModeHTML,
	}); err != nil {
		return err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	indexAlert := fmt.Sprintf("alert-notify-%s", time.Now().Format(time.DateOnly))
	res, err := repo.esClient.Index(
		indexAlert,
		bytes.NewReader(body),
		repo.esClient.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}
	return nil
}
