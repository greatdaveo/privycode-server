package models

import (
	"time"

	"gorm.io/gorm"
)

type ViewerLink struct {
	gorm.Model
	RepoName  string    `gorm:"not null" json:"repo_name"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"constraint:OnDelete:CASCADE"`
	Token     string    `gorm:"not null;unique" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	MaxViews  int       `json:"max_views"`
	ViewCount int       `json:"view_count"`
}
