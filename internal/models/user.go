package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email          string       `gorm:"unique;not null"`
	GitHubUsername string       `gorm:"unique;not null"`
	GitHubToken    string       `gorm:"not null"`
	ViewerLinks    []ViewerLink `gorm:"foreignKey:UserID"`
}
