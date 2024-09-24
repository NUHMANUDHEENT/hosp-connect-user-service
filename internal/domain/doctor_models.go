package domain

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	DoctorId        string
	Name            string
	Email           string `gorm:"uniqueIndex"`
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
	DoctorId int
	Name     string
}

type DoctorTokens struct {
	gorm.Model
	DoctorId     string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
