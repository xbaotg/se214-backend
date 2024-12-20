package models

import (
	"time"

	"github.com/google/uuid"
)

type Department struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	DepartmentName string
	DepartmentCode string
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
