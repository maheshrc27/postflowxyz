package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	config "github.com/maheshrc27/postflow/configs"
	"github.com/maheshrc27/postflow/internal/models"
	"github.com/maheshrc27/postflow/internal/repository"
	"github.com/maheshrc27/postflow/internal/transfer"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService interface {
	LoginCallback(ctx context.Context, code string) (err error, userID int64)
}

type authService struct {
	cfg config.Config
	u   repository.UserRepository
	c   repository.CreditsRepository
}

func NewAuthService(cfg config.Config, u repository.UserRepository, c repository.CreditsRepository) AuthService {
	return &authService{
		cfg: cfg,
		u:   u,
		c:   c,
	}
}

func (s *authService) LoginCallback(ctx context.Context, code string) (err error, userID int64) {

	if code == "" {
		err = errors.New("code or state is empty")
		slog.Info(err.Error())
		return err, 0
	}

	oauth2Config := &oauth2.Config{
		ClientID:     s.cfg.GoogleClientID,
		ClientSecret: s.cfg.GoogleClientSecret,
		RedirectURL:  "http://localhost:3000/login/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	if oauth2Config.ClientID == "" || oauth2Config.ClientSecret == "" || oauth2Config.RedirectURL == "" {
		err = errors.New("OAuth2 configuration is incomplete")
		slog.Info(err.Error())
		return err, 0
	}

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		slog.Info(err.Error())
		return err, 0
	}

	client := oauth2Config.Client(context.Background(), token)

	userInfo, err := GetUserInfo(client)
	if err != nil {
		return err, 0
	}

	user, isExist, err := s.u.GetByEmail(ctx, userInfo.Email)
	if err != nil {
		return err, 0
	}

	fmt.Printf("%v", user)

	if !isExist {
		userID, err = s.u.Create(ctx, &models.User{
			GoogleID:       userInfo.ID,
			Email:          userInfo.Email,
			Name:           userInfo.Name,
			ProfilePicture: userInfo.Picture,
		})

		credits := models.Credits{
			UserID:  userID,
			Credits: 1,
		}

		s.c.Create(ctx, &credits)

		if err != nil {
			return err, 0
		}
	} else {
		userID = user.ID
	}

	return nil, userID
}

func GetUserInfo(client *http.Client) (*transfer.GoogleUserInfo, error) {
	userInfoURL := "https://www.googleapis.com/oauth2/v1/userinfo"

	response, err := client.Get(userInfoURL)
	if err != nil {
		slog.Info(err.Error())
		return nil, fmt.Errorf("error fetching user info: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		slog.Info("Unexpected response status")
		return nil, fmt.Errorf("unexpected response status: %d", response.StatusCode)
	}

	var userInfo transfer.GoogleUserInfo
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		slog.Info(err.Error())
		return nil, fmt.Errorf("error decoding user info: %w", err)
	}

	return &userInfo, nil
}
