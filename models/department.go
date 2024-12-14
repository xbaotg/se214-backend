package models

import (
	"time"

	"github.com/google/uuid"
)

type Department struct {
	ID             uuid.UUID
	DepartmentName string
	DepartmentCode string
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
