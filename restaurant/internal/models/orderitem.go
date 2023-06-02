package models

import (
	"github.com/google/uuid"
	"time"
)

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	OrderUuid uuid.UUID `gorm:"type:uuid;not null"`
	Order     Order     `gorm:"foreignKey:OrderUuid"`
	CreatedAt time.Time `gorm:"default:current_timestamp;not null"`
}
