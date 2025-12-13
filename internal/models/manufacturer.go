package models

type Manufacturer struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null;unique" json:"name"`
	Founded     int    `json:"founded"`
	Defunct     *int   `json:"defunct,omitempty"`
	Country     string `json:"country"`
	Description string `gorm:"type:text" json:"description"`
}