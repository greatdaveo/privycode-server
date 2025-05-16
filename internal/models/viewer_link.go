package models

import (
	"time"

	"gorm.io/gorm"
)

type ViewerLink struct {
	gorm.Model
	RepoName  string `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	Token     string `gorm:"not null"`
	ExpiresAt time.Time
	MaxViews  int
	ViewCount int
}
