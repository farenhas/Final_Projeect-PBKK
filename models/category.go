package models

import "time"

type Category struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
