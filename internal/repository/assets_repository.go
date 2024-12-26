package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/maheshrc27/postflow/internal/models"
)

type MediaAssetRepository interface {
	Create(ctx context.Context, ma *models.MediaAsset) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.MediaAsset, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.MediaAsset, error)
}

type mediaAssetRepository struct {
	db *sql.DB
}

func NewMediaAssetRepository(db *sql.DB) MediaAssetRepository {
	return &mediaAssetRepository{db: db}
}

func (r *mediaAssetRepository) Create(ctx context.Context, ma *models.MediaAsset) (int64, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		slog.Info(err.Error())
		return 0, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO media_assets (user_id, file_name, file_type, file_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int64
	err = tx.QueryRowContext(ctx, query, ma.UserID, ma.FileName, ma.FileType, ma.FileURL).Scan(&id)
	if err != nil {
		slog.Info(err.Error())
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		slog.Info(err.Error())
		return 0, err
	}

	return id, nil
}

func (r *mediaAssetRepository) GetByID(ctx context.Context, id int64) (*models.MediaAsset, error) {
	query := `
		SELECT id, user_id, post_id, asset_url, created_at, updated_at
		FROM media_assets
		WHERE id = $1
	`

	var ma models.MediaAsset
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ma.ID,
		&ma.UserID,
		&ma.FileName,
		&ma.FileType,
		&ma.FileURL,
		&ma.ThumbnailURL,
		&ma.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Info(err.Error())
		return nil, err
	}

	return &ma, nil
}

func (r *mediaAssetRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.MediaAsset, error) {
	query := `SELECT id, file_name, file_type, file_url, created_at FROM media_assets WHERE user_id = $1`
	var rows *sql.Rows
	var err error

	rows, err = r.db.QueryContext(ctx, query, userID)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}
	defer rows.Close()

	var assets []*models.MediaAsset
	for rows.Next() {
		var asset models.MediaAsset
		err := rows.Scan(&asset.ID, &asset.FileName, &asset.FileType, &asset.FileURL, &asset.CreatedAt)
		if err != nil {
			slog.Info(err.Error())
			return nil, err
		}
		assets = append(assets, &asset)
	}
	return assets, nil
}
