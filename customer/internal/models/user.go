package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name       string    `gorm:"type:varchar(255);not null"`
	OfficeUuid uuid.UUID `gorm:"type:uuid;not null"`
	Office     Office    `gorm:"foreignKey:OfficeUuid"`
	CreatedAt  time.Time `gorm:"default:current_timestamp;not null"`
}
