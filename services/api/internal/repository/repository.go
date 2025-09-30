package repository

import (
	"context"
	"research-apm/services/api/internal/entity"
	"research-apm/services/api/internal/repository/internal/repository"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

type Repository interface {
	// get user
	GetUser(ctx context.Context) ([]entity.User, error)
	// create user
	CreateUser(ctx context.Context, data entity.User) error
	// get message
	GetMessage(ctx context.Context) ([]entity.Message, error)

	// get client do
	GetClientDO(ctx context.Context) ([]entity.ClientDo, error)
	// get profil
	GetProfil(ctx context.Context) ([]entity.Profil, error)
}

func NewRepository(mongoClient *mongo.Client, dbMessage *gorm.DB, dbClientDo *gorm.DB, dbProfil *gorm.DB, dbRedis *redis.Client) Repository {

	return repository.NewRepository(mongoClient, dbMessage, dbClientDo, dbProfil, dbRedis)
}
