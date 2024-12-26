package models

import "time"

type Credits struct {
	UserID    int64     `db:"ser_id" json:"ser_id"`
	Credits   int64     `db:"credits" json:"credits"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
