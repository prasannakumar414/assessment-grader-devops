package models

import "time"

const (
	StatusPending = "pending"
	StatusPassed  = "passed"
	StatusFailed  = "failed"
)

type Student struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Name          string     `json:"name" gorm:"not null"`
	Email         string     `json:"email" gorm:"not null;uniqueIndex"`
	RollNo        string     `json:"rollNo" gorm:"not null;uniqueIndex"`
	Status        string     `json:"status" gorm:"not null;default:pending"`
	LastCheckedAt *time.Time `json:"lastCheckedAt"`
	ErrorMessage  string     `json:"errorMessage"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
