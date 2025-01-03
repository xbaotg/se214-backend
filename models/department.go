package models

import (
	"time"

	"github.com/google/uuid"
)

type Department struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	DepartmentName string
	DepartmentCode string `gorm:"unique;not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
