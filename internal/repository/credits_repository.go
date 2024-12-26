// repository/user_repository.go
package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/maheshrc27/postflow/internal/models"
)

type CreditsRepository interface {
	GetByUserID(ctx context.Context, id int64) (*models.Credits, bool, error)
	Create(ctx context.Context, credits *models.Credits) (int64, error)
	UpdateCredits(ctx context.Context, credits, userID int64) error
}

type creditsRepository struct {
	db *sql.DB
}

func NewCreditsRepository(db *sql.DB) CreditsRepository {
	return &creditsRepository{db: db}
}

func (r *creditsRepository) GetByUserID(ctx context.Context, id int64) (*models.Credits, bool, error) {
	var credits models.Credits
	query := "SELECT user_id, credits FROM credits WHERE user_id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&credits.UserID, &credits.Credits)
	if err != nil {
		slog.Info(err.Error())
		return nil, false, err
	}
	return &credits, true, nil
}

func (r *creditsRepository) Create(ctx context.Context, credits *models.Credits) (int64, error) {
	query := "INSERT INTO credits (user_id, credits) VALUES ($1, $2) RETURNING user_id"
	var id int64
	err := r.db.QueryRowContext(ctx, query, credits.UserID, credits.Credits).Scan(&id)
	if err != nil {
		slog.Info(err.Error())
		return 0, err
	}
	return id, nil
}

func (r *creditsRepository) UpdateCredits(ctx context.Context, credits, userID int64) error {
	query := `
		UPDATE credits 
		SET credits = $1,
			updated_at = $2
		WHERE user_id = $3
	`
	_, err := r.db.ExecContext(ctx, query, credits, time.Now(), userID)
	if err != nil {
		slog.Info(err.Error())
		return err
	}

	return nil
}
