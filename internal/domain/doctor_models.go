package domain

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	DoctorId        string
	Name            string
	Email           string `gorm:"unique"`
	Password        string
	Phone           int
	SpecilazationId int
	AvailabilityId  int
	Role            string
}
type DoctorAvailability struct {
	gorm.Model
	DoctorId  int
	Day       string
	StartTime string
	EndTime   string
	Available bool
}
type DoctorSpecialization struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Description string
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
	DoctorID  string    `json:"doctor_id"`
	EventType string    `json:"day_of_week"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
