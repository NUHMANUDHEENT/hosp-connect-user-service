package service

import (
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/repository"
)

type PatientService interface {
	SignIn(patient domain.Patient) (string, error)
	SignUp(patient domain.Patient) (string, error)
	GetProfile(patientId string) (domain.Patient, error)
	UpdateProfile(patient domain.Patient) error
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
	return resp, nil
}

func (p *patientService) GetProfile(patientId string) (domain.Patient, error) {
    return p.repo.GetProfile(patientId)
}

func (p *patientService) UpdateProfile(patient domain.Patient) error {
    return p.repo.UpdateProfile(patient)
}