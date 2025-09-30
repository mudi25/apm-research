package config

import (
	"context"
	"fmt"
	"os"
	"research-apm/services/alert/internal/repository"
	"research-apm/services/alert/internal/service"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-telegram/bot"
)

type App struct {
	schduler gocron.Scheduler
}

func (a *App) Shutdown() {
	if err := a.schduler.Shutdown(); err != nil {
		fmt.Println("[ERROR] shutdown schduler", err.Error())
	}
}

func NewApp(ctx context.Context) (*App, error) {
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: strings.Split(os.Getenv("ELASTIC_HOST"), ","),
		Username:  os.Getenv("ELASTIC_USER"),
		Password:  os.Getenv("ELASTIC_PASS"),
	})
	if err != nil {
		return nil, err
	}
	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return nil, err
	}
	scheduler, err := gocron.NewScheduler()
	repo := repository.NewRepository(
		esClient,
		b,
		os.Getenv("TELEGRAM_CHAT_ID"),
	)
	scheduler.NewJob(
		gocron.CronJob("0 * * * * *", true),
		gocron.NewTask(service.SendAlert, ctx, repo),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		return nil, err
	}
	scheduler.Start()
	return &App{schduler: scheduler}, nil

}
