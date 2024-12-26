package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/maheshrc27/postflow/internal/repository"
)

type CreditsService interface {
	GetCredits(ctx context.Context, id int64) (int64, error)
}

type creditsService struct {
	c repository.CreditsRepository
}

func NewCreditsService(c repository.CreditsRepository) CreditsService {
	return &creditsService{
		c: c,
	}
}

func (s *creditsService) GetCredits(ctx context.Context, id int64) (int64, error) {
	credits, isExist, err := s.c.GetByUserID(ctx, id)
	if err != nil {
		return 0, err
	}

	if !isExist {
		err = errors.New("User not found")
		slog.Info(err.Error())
	}

	return credits.Credits, nil
}

func (s *creditsService) UpdateCredits(ctx context.Context, userID, newCredits int64) error {
	err := s.c.UpdateCredits(ctx, newCredits, userID)
	if err != nil {
		return err
	}

	return nil
}
