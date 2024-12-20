package models

import (
	"time"
)

type AllCourses struct {
	CourseName string    `gorm:"type:text;primaryKey"`
	Status    bool	`gorm:"type:boolean;default:true"`
	CourseFullname string `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	// Course    Course  `gorm:"foreignKey:CourseName;references:CourseName" json:"-"`
}