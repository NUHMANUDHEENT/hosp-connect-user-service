package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PatientRepository interface {
	SignIn(patient domain.Patient) (string, error)
	SignUp(patient domain.Patient) (string, error)
	SignUpVerify(email string) (string, error)
	DeletePatient(patientID string) (string, error)
	Block(patientID string, reason string) (string, error)
	GetProfile(patientId string) (domain.Patient, error)
	UpdateProfile(patient domain.Patient) error
	ListPatients() ([]domain.Patient, error)
}
type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

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

	return patient.PatientID, nil
}

// Patient SignUp Repository logic
func (p *patientRepository) SignUp(patient domain.Patient) (string, error) {
	var existingPatient domain.Patient
	if err := p.db.First(&existingPatient, "email = ?", patient.Email).Error; err == nil {
		return "", errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(patient.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	uuid := uuid.New()
	existingPatient.PatientID = fmt.Sprintf("pt-%s", uuid.String())

	patient.Password = string(hashedPassword)
	patient.VerifyStatus = false
	fmt.Println(existingPatient)

	// Save the new patient to the database
	if err := p.db.Create(&patient).Error; err != nil {
		return "", err
	}

	return "Successfully registered", nil
}
func (p *patientRepository) SignUpVerify(email string) (string, error) {
	var patient domain.Patient
	if err := p.db.First(&patient, "email = ?", email).Error; err != nil {
		return "", err
	}
	patient.VerifyStatus = true
	if err := p.db.Where("email = ?", patient.Email).Model(&patient).Updates(patient).Error; err != nil {
		return "email not found ", err
	}
	return "Email has been successfully verified!" + email, nil
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
	if err := p.db.Where("email = ?", patient.Email).Model(&patient).Updates(patient).Error; err != nil {
		return "blocking failed", err
	}
	return "Patient blocked successfully", nil
}

func (p *patientRepository) GetProfile(patientId string) (domain.Patient, error) {
	var patient domain.Patient
	if err := p.db.Where("patient_id = ?", patientId).First(&patient).Error; err != nil {
		return patient, err
	}
	return patient, nil
}

func (p *patientRepository) UpdateProfile(patient domain.Patient) error {
	if err := p.db.Where("patient_id = ?", patient.Email).Model(&patient).Updates(patient).Error; err != nil {
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
