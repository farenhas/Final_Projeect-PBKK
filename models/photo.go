package models

import "time"

type Photo struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	CategoryID  uint      `gorm:"not null"`
	ImagePath   string    `gorm:"size:255;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	Category    Category  `gorm:"foreignKey:CategoryID"`
}
