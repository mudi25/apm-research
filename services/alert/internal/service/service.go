package service

import (
	"context"
	"fmt"
	"research-apm/services/alert/internal/repository"

	"golang.org/x/sync/errgroup"
)

func SendAlert(ctx context.Context, repo repository.Repository) error {
	items, err := repo.GetAlert(ctx)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(20)
	for _, it := range items {
		fmt.Println("[INFO] new alert", it.ServiceName)
		g.Go(func() error {
			if err := repo.SendTelegram(ctx, it); err != nil {
				fmt.Println(err.Error())
			}
			return nil
		})
	}
	return g.Wait()
}
