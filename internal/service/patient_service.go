package service

import (
	"github.com/go-redis/redis"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/config"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/repository"
)

type PatientService interface {
	SignIn(patient domain.Patient) (string, error)
	SignUp(patient domain.Patient) (string, error)
	SignUpVerify(token string) (string, error)
	GetProfile(patientId string) (domain.Patient, error)
	UpdateProfile(patient domain.Patient) error
	AddPrescription(data domain.PatientPrescription) error
	GetPrescription(patientId, query string) ([]domain.PatientPrescription, error)
}

type patientService struct {
	repo repository.PatientRepository
}

func NewPatientService(repo repository.PatientRepository) PatientService {
	return &patientService{repo: repo}
}

// SignIn logic for patient
func (p *patientService) SignIn(patient domain.Patient) (string, error) {
	resp, err := p.repo.SignIn(patient)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// SignUp logic for patient
func (p *patientService) SignUp(patient domain.Patient) (string, error) {
	resp, err := p.repo.SignUp(patient)
	if err != nil {
		return resp, err
	}
	resp, err = config.SignUpverify(patient.Email)
	if err != nil {
		return "failed to send verification to mail", err
	}

	return resp, nil
}
func (p *patientService) SignUpVerify(token string) (string, error) {
	email, err := config.Rdb.Get(token).Result()
	if err == redis.Nil {
		return "Token expired or invalid", err
	} else if err != nil {
		return "Failed to verify token", err
	}
	resp, err := p.repo.SignUpVerify(email)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (p *patientService) GetProfile(patientId string) (domain.Patient, error) {
	return p.repo.GetProfile(patientId)
}

func (p *patientService) UpdateProfile(patient domain.Patient) error {
	return p.repo.UpdateProfile(patient)
}

func (p *patientService) AddPrescription(data domain.PatientPrescription) error {
	return p.repo.StorePatientPrescription(data)
}
func (p patientService) GetPrescription(patientId, query string) ([]domain.PatientPrescription, error) {
	return p.repo.GetPrescriptions(patientId, query)
}
