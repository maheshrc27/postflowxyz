package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/maheshrc27/postflow/internal/models"
	"github.com/maheshrc27/postflow/internal/repository"
)

type UserService interface {
	GetUserInfo(ctx context.Context, id int64) (*models.User, error)
	RemoveUser(ctx context.Context, userID int64) error
}

type userService struct {
	u repository.UserRepository
}

func NewUserService(u repository.UserRepository) UserService {
	return &userService{
		u: u,
	}
}

func (s *userService) GetUserInfo(ctx context.Context, id int64) (*models.User, error) {
	user, isExist, err := s.u.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !isExist {
		err = errors.New("User not found")
		slog.Info(err.Error())
	}

	return user, nil
}

func (s *userService) RemoveUser(ctx context.Context, userID int64) error {
	err := s.u.Remove(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
