package models

import "time"

type MediaAsset struct {
	ID           int64     `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	FileName     string    `db:"file_name" json:"file_name"`
	FileType     string    `db:"file_type" json:"file_type"`
	FileURL      string    `db:"file_url" json:"file_url"`
	ThumbnailURL string    `db:"thumbnail_url" json:"thumbnail_url"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
