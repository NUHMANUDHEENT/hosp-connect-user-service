package repository

import (
	"errors"

	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PatientRepository interface {
	SignIn(patient domain.Patient) (string, error)
	SignUp(patient domain.Patient) (string, error)
	// New methods for patient management
	DeletePatient(patientID string) (string, error)
	Block(patientID string, reason string) (string, error)
	GetProfile(patientId int32) (domain.Patient, error)
	UpdateProfile(patient domain.Patient) error
	ListPatients() ([]domain.Patient, error)
}
type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

// Patient SignIn Repository logic
func (p *patientRepository) SignIn(patient domain.Patient) (string, error) {
	var patientCheck domain.Patient
	if err := p.db.First(&patientCheck, "email = ?", patient.Email).Error; err != nil {
		return "Email or password is incorrect", err
	}
	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(patientCheck.Password), []byte(patient.Password)); err != nil {
		return "Email or password is incorrect", err
	}
	if patient.IsBlock {
		return "Your account is blocked", errors.New("account is blocked")
	}

	return "Successfully logged in", nil
}

// Patient SignUp Repository logic
func (p *patientRepository) SignUp(patient domain.Patient) (string, error) {
	// Check if email already exists
	var existingPatient domain.Patient
	if err := p.db.First(&existingPatient, "email = ?", patient.Email).Error; err == nil {
		return "", errors.New("email already in use")
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(patient.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	patient.Password = string(hashedPassword)

	// Save the new patient to the database
	if err := p.db.Create(&patient).Error; err != nil {
		return "", err
	}

	return "Successfully registered", nil
}
func (p *patientRepository) DeletePatient(patientID string) (string, error) {
	var patient domain.Patient
	if err := p.db.First(&patient, patientID).Error; err != nil {
		return "Patient not found", err
	}
	if err := p.db.Delete(&patient).Error; err != nil {
		return "Patient deletion failed", err
	}
	return "Patient deleted successfully", nil
}

// Block logic for patient
func (p *patientRepository) Block(patientID string, reason string) (string, error) {
	var patient domain.Patient
	if err := p.db.First(&patient, patientID).Error; err != nil {
		return "Patient not found", err
	}
	patient.IsBlock = true
	patient.IsBlockReason = reason
	if err := p.db.Save(&patient).Error; err != nil {
		return "Patient blocking failed", err
	}
	return "Patient blocked successfully", nil
}

func (p *patientRepository) GetProfile(patientId int32) (domain.Patient, error) {
	var patient domain.Patient
	if err := p.db.Where("parient_id = ?", patientId).First(&patient).Error; err != nil {
		return patient, err
	}
	return patient, nil
}

func (p *patientRepository) UpdateProfile(patient domain.Patient) error {
	if err := p.db.Where("email = ?", patient.Email).Model(&patient).Updates(patient).Error; err != nil {
		return err
	}
	return nil
}

func (p *patientRepository) ListPatients() ([]domain.Patient, error) {
	var patients []domain.Patient
	if err := p.db.Find(&patients).Error; err != nil {
		return nil, err
	}
	return patients, nil
}
