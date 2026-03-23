package models

import "gorm.io/gorm"

type Contest struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Code     string `gorm:"uniqueIndex:idx_platform_code"`
	Platform string `gorm:"not null;uniqueIndex:idx_platform_code"`
	Start    string
	End      string
}
