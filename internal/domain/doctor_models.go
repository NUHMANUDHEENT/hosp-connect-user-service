package domain

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	DoctorId         string
	Name             string
	Email            string `gorm:"unique"`
	Password         string
	Phone            int
	SpecializationId int
	AvailabilityId   int
	Role             string
}

type DoctorTokens struct {
	gorm.Model
	DoctorId     string `gorm:"unique"`
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
type AvailabilitySlot struct {
	gorm.Model
	DoctorID   string `json:"doctor_id"`
	DoctorName string
	EventType  string    `json:"day_of_week"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}
type AvailableDates struct{
	DateTime time.Time
	IsAvailable string
}