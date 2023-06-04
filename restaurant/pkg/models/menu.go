package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type Menu struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	OnDate          time.Time      `gorm:"not null"`
	OpeningRecordAt time.Time      `gorm:"not null"`
	ClosingRecordAt time.Time      `gorm:"not null"`
	Year            int            `gorm:"type:int;not null"`
	Month           int            `gorm:"type:int;not null"`
	Day             int            `gorm:"type:int;not null"`
	ProductsUuids   pq.StringArray `gorm:"type:text[]"`
	CreatedAt       time.Time      `gorm:"default:current_timestamp;not null"`
}
