package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lahaehae/crud_project/internal/models"
	"github.com/lahaehae/crud_project/internal/repository"
)

type SubUserService struct {
    repo *repository.SubUserRepository
}

func NewSubUserService(repo *repository.SubUserRepository) *SubUserService {
    return &SubUserService{
        repo: repo,
    }
}

func (s *SubUserService) CreateSubUser(ctx context.Context, ownerID string, name, email string) (*models.SubUser, error) {
    now := time.Now()
    subUser := &models.SubUser{
        ID:        uuid.New().String(),
        OwnerID:   ownerID,
        Name:      name,
        Email:     email,
        CreatedAt: now,
        UpdatedAt: now,
    }

    if err := s.repo.Create(ctx, subUser); err != nil {
        return nil, err
    }
    return subUser, nil
}

func (s *SubUserService) GetUserSubUsers(ctx context.Context, ownerID string) ([]models.SubUser, error) {
    return s.repo.GetByOwner(ctx, ownerID)
}

// Другие методы сервиса...