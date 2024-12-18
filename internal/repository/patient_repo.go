package repository

import (
	"errors"
	"fmt"
	"time"

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
	StorePatientPrescription(data domain.PatientPrescription) error
	GetPrescriptions(patientId, query string) ([]domain.PatientPrescription, error)
	GetPatientCount() (int, error)
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

	return patientCheck.PatientID, nil
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
	patient.PatientID = fmt.Sprintf("pt-%s", uuid.String())

	patient.Password = string(hashedPassword)
	patient.VerifyStatus = false

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
	if err := p.db.First(&patient, "patient_id", patientID).Error; err != nil {
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
	if err := p.db.First(&patient, "patient_id", patientID).Error; err != nil {
		return "Patient not found", err
	}
	if patient.IsBlock {
		patient.IsBlock = false
		patient.IsBlockReason = ""
	} else {
		patient.IsBlock = true
		patient.IsBlockReason = reason
	}
	if err := p.db.Where("id = ?", patient.ID).Model(&patient).Updates(patient).Error; err != nil {
		return "blocking failed", err
	}
	if !patient.IsBlock {
		return "Patient Unblocked successfully", nil
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
	if err := p.db.Where("patient_id = ?", patient.PatientID).Model(&patient).Updates(patient).Error; err != nil {
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
func (p *patientRepository) StorePatientPrescription(data domain.PatientPrescription) error {

	var patient domain.Patient
	if err := p.db.First(&patient, "patient_id =?", data.PatientId).Error; err != nil {
		return errors.New("patient not found")
	}
	if err := p.db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}
func (p *patientRepository) GetPrescriptions(patientId, query string) ([]domain.PatientPrescription, error) {
	var prescriptions []domain.PatientPrescription
	currentDate := time.Now()

	switch query {
	case "today":
		if err := p.db.Where("patient_id = ? AND DATE(created_at) = ?", patientId, currentDate.Format("2006-01-02")).Find(&prescriptions).Error; err != nil {
			return nil, err
		}
	case "old":
		if err := p.db.Where("patient_id = ? AND DATE(created_at) < ?", patientId, currentDate.Format("2006-01-02")).Find(&prescriptions).Error; err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid query type")
	}
	return prescriptions, nil
}
func (p *patientRepository) GetPatientCount() (int, error) {
	var count int64
	err := p.db.Model(&domain.Patient{}).Count(&count).Error
	return int(count), err
}
