package models

import (
	"time"

	"github.com/google/uuid"
)

type PrerequisiteCourse struct {
	CourseID       uuid.UUID
	PrerequisiteID uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Course         Course `gorm:"foreignKey:CourseID;references:ID"`
	Prerequisite   Course `gorm:"foreignKey:PrerequisiteID;references:ID"`
}
