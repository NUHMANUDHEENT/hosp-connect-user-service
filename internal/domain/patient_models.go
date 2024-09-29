package domain

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	PatientID  string   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Password string
	Name     string
	Phone    int
	Age      int32
	Gender   string
	VerifyStatus bool
	IsBlock   bool
	IsBlockReason string
}
