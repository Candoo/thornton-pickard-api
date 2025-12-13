package models

import (
	"time"

	"gorm.io/gorm"
)

type Ephemera struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Type           string    `gorm:"not null" json:"type"` // catalog, manual, advertisement, etc.
	Title          string    `gorm:"not null" json:"title"`
	Year           int       `json:"year"`
	Pages          *int      `json:"pages,omitempty"`
	Description    string    `gorm:"type:text" json:"description"`
	ScanURL        string    `json:"scan_url"`
	ThumbnailURL   string    `json:"thumbnail_url"`
	RelatedCameras string    `json:"related_cameras"` // Store as JSON array of camera IDs
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}