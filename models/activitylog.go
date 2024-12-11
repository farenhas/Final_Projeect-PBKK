package models

import "time"

type ActivityLog struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `json:"user_id"`
    Action    string    `json:"action"` // CRUD: Create, Read, Update, Delete
    Entity    string    `json:"entity"` // Contoh: Photo, Category
    Timestamp time.Time `gorm:"autoCreateTime"`
}
