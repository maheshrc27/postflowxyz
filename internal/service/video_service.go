package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/maheshrc27/postflow/internal/models"
	"github.com/maheshrc27/postflow/internal/repository"
	"github.com/maheshrc27/postflow/internal/transfer"
)

type VideoService interface {
	GetVideos(ctx context.Context, userID int64) ([]*models.MediaAsset, error)
	RequestVideo(ctx context.Context, userID int64, jsonData string) (string, error)
}

type videoService struct {
	c repository.CreditsRepository
	a repository.MediaAssetRepository
}

func NewVideoService(c repository.CreditsRepository, a repository.MediaAssetRepository) VideoService {
	return &videoService{
		c: c,
		a: a,
	}
}

func (s *videoService) GetVideos(ctx context.Context, userID int64) ([]*models.MediaAsset, error) {
	videos, err := s.a.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (s *videoService) RequestVideo(ctx context.Context, userID int64, jsonData string) (string, error) {
	credits, isExist, err := s.c.GetByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

	if !isExist {
		err = errors.New("Credits doesn't exist")
		slog.Info(err.Error())
		return "", err
	}

	if credits.Credits < 1 {
		err = errors.New("Not enought credits")
		slog.Info(err.Error())
		return "", err
	}

	url := "http://127.0.0.1:5000/generate"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		slog.Info(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New("Non-OK HTTP status")
		slog.Info(err.Error())
		return "", err
	}

	var response transfer.VideoResponseTransfer
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		slog.Info(err.Error())
		return "", err
	}

	oldCredits := credits.Credits
	newCredits := oldCredits - 1

	err = s.c.UpdateCredits(ctx, newCredits, userID)
	if err != nil {
		return "", err
	}

	videoURL := fmt.Sprintf("https://pub-f8f43aa198a449518df6744ec9ce452c.r2.dev/%s.mp4", response.VideoID)

	asset := models.MediaAsset{
		UserID:   userID,
		FileName: response.VideoID,
		FileType: "video/mp4",
		FileURL:  videoURL,
	}

	_, err = s.a.Create(ctx, &asset)
	if err != nil {
		return "", err
	}

	return videoURL, nil
}
