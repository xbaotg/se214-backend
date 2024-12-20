package models

import (
	"time"
)

type PrerequisiteCourse struct {
	CourseID       string    `gorm:"type:text;not null"`
	PrerequisiteID string    `gorm:"type:text;not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

	Course       AllCourses `gorm:"foreignKey:CourseID;references:CourseName" json:"-"`
	Prerequisite AllCourses `gorm:"foreignKey:PrerequisiteID;references:CourseName" json:"-"`
}