package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ShortLink struct {
	gorm.Model
	OriginalURL string     `gorm:"type:text;not null"`
	Code        string     `gorm:"type:varchar(10);unique_index;not null"`
	ExpiresAt   *time.Time `gorm:"index;default:NULL"`
	ViewCount   uint       `gorm:"default:0"`
	IsDeleted   bool       `gorm:"default:false"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&ShortLink{}).Error
}
