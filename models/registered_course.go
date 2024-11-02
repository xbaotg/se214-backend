package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CoStatus string

const (
	CoStatusDone        CoStatus = "done"
	CoStatusFailed      CoStatus = "failed"
	CoStatusProgressing CoStatus = "progressing"
)

func (e *CoStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CoStatus(s)
	case string:
		*e = CoStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for CoStatus: %T", src)
	}
	return nil
}

type NullCoStatus struct {
	CoStatus CoStatus
	Valid    bool // Valid is true if CoStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCoStatus) Scan(value interface{}) error {
	if value == nil {
		ns.CoStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CoStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCoStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.CoStatus), nil
}

type RegisteredCourse struct {
	CourseID       uuid.UUID
	UserID         uuid.UUID
	CourseYear     int32
	CourseSemester int32
	Status         CoStatus
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

	Course Course `gorm:"foreignKey:CourseID;references:ID"`
	User   User   `gorm:"foreignKey:UserID"`
}
