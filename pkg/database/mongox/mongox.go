package mongox

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	Uri             string
	MinPoolSize     uint64
	MaxPoolSize     uint64
	MaxConnIdleTime time.Duration
}

func NewClient(cfg Config) (*mongo.Client, error) {
	clientOpts := options.Client().
		ApplyURI(cfg.Uri).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMaxConnIdleTime(cfg.MaxConnIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("connect failed: %s", err.Error())
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping failed: %s", err.Error())
	}

	return client, nil
}

func Disconnect(dbClient *mongo.Client) {
	if dbClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := dbClient.Disconnect(ctx); err != nil {
			fmt.Println("failed shutdown mongodb", err.Error())
		}
	}
}
