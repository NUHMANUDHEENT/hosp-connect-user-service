package domain

import (
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	PatientID     string
	Email         string `gorm:"unique"`
	Password      string
	Name          string
	Phone         int
	Age           int32
	Gender        string
	VerifyStatus  bool
	IsBlock       bool
	IsBlockReason string
}
type PatientPrescription struct {
	gorm.Model
	PatientId    string `json:"patient_id"`
	DoctorId     string `json:"doctor_id"`
	Prescription string `json:"prescription" gorm:"type:jsonb"`  // Store as JSONB in PostgreSQL
}

type Prescription struct {
	Medication string `json:"medication"`
	Dosage     string `json:"dosage"`
	Frequency  string `json:"frequency"`
}