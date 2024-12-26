package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	config "github.com/maheshrc27/postflow/configs"
	"github.com/maheshrc27/postflow/internal/models"
	"github.com/maheshrc27/postflow/internal/repository"
)

const (
	productId      = "ehajql"
	price1         = 500
	price2         = 900
	price3         = 1500
	InitialCredits = 1
	CreditsPrice1  = 10
	CreditsPrice2  = 30
	CreditsPrice3  = 50
)

type PaymentService interface {
	HandlePayment(ctx context.Context, email, productID, productPrice string) error
}

type paymentService struct {
	cfg config.Config
	u   repository.UserRepository
	c   repository.CreditsRepository
}

func NewPaymentService(cfg config.Config, u repository.UserRepository, c repository.CreditsRepository) PaymentService {
	return &paymentService{
		cfg: cfg,
		u:   u,
		c:   c,
	}
}

func (s *paymentService) HandlePayment(ctx context.Context, email, productID, productPrice string) error {

	if productID != productId {
		err := errors.New("No product exists like this")
		slog.Info(err.Error())
		return fmt.Errorf("fetching user by email failed: %w", err)
	}

	price, _ := strconv.Atoi(productPrice)

	user, isExist, err := s.u.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("fetching user by email failed: %w", err)
	}

	var userID int64
	if !isExist {
		userID, err = s.createUserAndCredits(ctx, email)
		if err != nil {
			return err
		}
	} else {
		userID = user.ID
	}

	credits, isExist, err := s.c.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("fetching credits for user %d failed: %w", userID, err)
	}

	if !isExist {
		err = fmt.Errorf("credits row for user_id %d doesn't exist", userID)
		slog.Info(err.Error())
		return err
	}

	var newCredits int64
	switch price {
	case price1:
		newCredits = credits.Credits + CreditsPrice1
	case price2:
		newCredits = credits.Credits + CreditsPrice2
	case price3:
		newCredits = credits.Credits + CreditsPrice3
	default:
		return fmt.Errorf("invalid productID: %s", productID)
	}

	if err := s.c.UpdateCredits(ctx, newCredits, userID); err != nil {
		slog.Error("failed to update credits", "error", err, "userID", userID)
		return fmt.Errorf("updating credits failed: %w", err)
	}

	return nil
}

func (s *paymentService) createUserAndCredits(ctx context.Context, email string) (int64, error) {
	newUser := models.User{Email: email}
	userID, err := s.u.Create(ctx, &newUser)
	if err != nil {
		return 0, fmt.Errorf("creating user failed: %w", err)
	}

	newCredits := models.Credits{
		UserID:  userID,
		Credits: InitialCredits,
	}
	if _, err = s.c.Create(ctx, &newCredits); err != nil {
		return 0, fmt.Errorf("creating credits row failed for user %d: %w", userID, err)
	}

	return userID, nil
}
