package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Day string

const (
	DayMonday    Day = "monday"
	DayTuesday   Day = "tuesday"
	DayWednesday Day = "wednesday"
	DayThursday  Day = "thursday"
	DayFriday    Day = "friday"
	DaySaturday  Day = "saturday"
	DaySunday    Day = "sunday"
)

func (e *Day) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Day(s)
	case string:
		*e = Day(s)
	default:
		return fmt.Errorf("unsupported scan type for Day: %T", src)
	}
	return nil
}

type NullDay struct {
	Day   Day
	Valid bool // Valid is true if Day is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullDay) Scan(value interface{}) error {
	if value == nil {
		ns.Day, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Day.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullDay) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Day), nil
}

type Course struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CourseTeacherID uuid.UUID `gorm:"type:uuid;not null"`
	DepartmentID    uuid.UUID `gorm:"type:uuid;not null"`
	CourseName      string    `gorm:"type:text;not null"`
	CourseFullname  string    `gorm:"type:text;not null"`
	CourseCredit    int32       `gorm:"not null"`
	CourseYear      int32       `gorm:"not null"`
	CourseSemester  int32       `gorm:"not null"`
	CourseStartShift int32      `gorm:"not null"`
	CourseEndShift   int32      `gorm:"not null"`
	CourseDay       Day       `gorm:"type:day;not null"`
	Confirmed       bool      `gorm:"not null;default:false"`
	MaxEnroller     int32       `gorm:"not null"`
	CurrentEnroller int32      `gorm:"not null"`
	CourseRoom      string    `gorm:"type:text;not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`

	Teacher    User       `gorm:"foreignKey:CourseTeacherID" json:"-"`
	Department Department `gorm:"foreignKey:DepartmentID" json:"-"`
}
