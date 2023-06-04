package models

import (
	"github.com/comp1x/final-task/customer/pkg/models"
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID          uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ProductUuid uuid.UUID   `gorm:"type:uuid;not null"`
	Product     Product     `gorm:"foreignKey:ProductUuid"`
	UserUuid    uuid.UUID   `gorm:"type:uuid;not null"`
	User        models.User `gorm:"foreignKey:UserUuid"`
	Count       int64       `gorm:"type:int; not null"`
	CreatedAt   time.Time   `gorm:"default:current_timestamp;not null"`
}
