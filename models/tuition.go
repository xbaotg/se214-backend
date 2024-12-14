package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TuStatus string

const (
	TuStatusPaid   TuStatus = "paid"
	TuStatusUnpaid TuStatus = "unpaid"
)

func (e *TuStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TuStatus(s)
	case string:
		*e = TuStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for TuStatus: %T", src)
	}
	return nil
}

type NullTuStatus struct {
	TuStatus TuStatus
	Valid    bool // Valid is true if TuStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTuStatus) Scan(value interface{}) error {
	if value == nil {
		ns.TuStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TuStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTuStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TuStatus), nil
}

type Tuition struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Tuition         int32
	Paid            int32
	TotalCredit     int32
	Year            int32
	Semester        int32
	TuitionStatus   TuStatus
	TuitionDeadline time.Time
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
