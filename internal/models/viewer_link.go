package models

import (
	"time"

	"gorm.io/gorm"
)

type ViewerLink struct {
	gorm.Model
	RepoName  string `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	User      User   `gorm:"constraint:OnDelete:CASCADE"`
	Token     string `gorm:"not null;unique"`
	ExpiresAt time.Time
	MaxViews  int
	ViewCount int
}
