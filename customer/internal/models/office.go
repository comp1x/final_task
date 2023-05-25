package models

import (
	"github.com/google/uuid"
	"time"
)

type Office struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Address   string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"default:current_timestamp;not null"`
}
