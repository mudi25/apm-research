package service

import (
	"context"
	"research-apm/pkg/errors"
	"research-apm/pkg/errors/codes"
	"research-apm/pkg/tracer"
	"research-apm/services/api/internal/entity"
	"research-apm/services/api/internal/repository"
	"time"

	"github.com/oklog/ulid/v2"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

// get user
func (service *Service) GetUser(ctx context.Context) ([]entity.User, error) {
	ctx, span := tracer.StartSpan(ctx, "service.GetUser")
	defer span.End()
	result, err := service.repo.GetUser(ctx)
	if err != nil {
		return nil, errors.Wrap(codes.Internal, "gagal mencari data user", err)
	}
	return result, nil
}

// create user
func (service *Service) CreateUser(ctx context.Context, data entity.User) (string, error) {
	ctx, span := tracer.StartSpan(ctx, "service.CreateUser")
	defer span.End()
	data.ID = ulid.Make().String()
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	if err := service.repo.CreateUser(ctx, data); err != nil {
		return "", errors.Wrap(codes.Internal, "gagal membuat user", err)
	}
	return data.ID, nil

}

// get message
func (service *Service) GetMessage(ctx context.Context) ([]entity.Message, error) {
	ctx, span := tracer.StartSpan(ctx, "service.GetMessage")
	defer span.End()
	result, err := service.repo.GetMessage(ctx)
	if err != nil {
		return nil, errors.Wrap(codes.Internal, "gagal mencari data message", err)
	}
	return result, nil
}

// get client do
func (service *Service) GetClientDO(ctx context.Context) ([]entity.ClientDo, error) {
	ctx, span := tracer.StartSpan(ctx, "service.GetClientDO")
	defer span.End()
	result, err := service.repo.GetClientDO(ctx)
	if err != nil {
		tracer.CaptureError(ctx, err)
		return nil, errors.Wrap(codes.Internal, "gagal mencari data client do", err)
	}
	return result, nil
}

// get profil
func (service *Service) GetProfil(ctx context.Context) ([]entity.Profil, error) {
	ctx, span := tracer.StartSpan(ctx, "service.GetProfil")
	defer span.End()
	result, err := service.repo.GetProfil(ctx)
	if err != nil {
		return nil, errors.Wrap(codes.Internal, "gagal mencari data profil", err)
	}
	return result, nil
}
