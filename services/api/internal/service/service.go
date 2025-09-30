package service

import (
	"context"
	"research-apm/services/api/internal/entity"
	"research-apm/services/api/internal/repository"
	"research-apm/services/api/internal/service/internal/service"
)

type Service interface {
	// get user
	GetUser(ctx context.Context) ([]entity.User, error)

	// create user
	CreateUser(ctx context.Context, data entity.User) (string, error)

	// get message
	GetMessage(ctx context.Context) ([]entity.Message, error)

	// get client do
	GetClientDO(ctx context.Context) ([]entity.ClientDo, error)

	// get profil
	GetProfil(ctx context.Context) ([]entity.Profil, error)
}

func NewService(repo repository.Repository) Service {
	return service.NewService(repo)
}
