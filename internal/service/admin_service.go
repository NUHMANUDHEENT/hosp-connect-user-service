package service

import (
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/repository"
)

type AdminService interface {
	SignIn(admin domain.Admin) (string, error)
	AddDoctor(email, name, password string, phone, specializationId int32) (string, error)
	DeleteDoctor(doctorId string) (string,error)
	// New methods for managing patients
	AddPatient(patient domain.Patient) (string, error)
	DeletePatient(patientID string) (string, error)
	BlockPatient(patientID string, reason string) (string, error)
	ListDoctors() ([]domain.Doctor, error)
    ListPatients() ([]domain.Patient, error)
}
type adminService struct {
	repo       repository.AdminRepository
	doctorRepo repository.DoctorRepository
	patientRepo repository.PatientRepository // Inject the PatientRepository
}

// NewAdminService constructor
func NewAdminService(repo repository.AdminRepository, docrepo repository.DoctorRepository, patientRepo repository.PatientRepository) AdminService {
	return &adminService{repo: repo, doctorRepo: docrepo, patientRepo: patientRepo}
}

func (a *adminService) SignIn(admin domain.Admin) (string, error) {
	resp, err := a.repo.SignIn(admin)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (a *adminService) AddDoctor(email, name, password string, phone, specializationId int32) (string, error) {
	resp, err := a.doctorRepo.SignUpStore(domain.Doctor{
		Email:           email,
		Name:            name,
		Password:        password,
		SpecilazationId: int(specializationId),
		Phone:           int(phone),
	})
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (a *adminService) DeleteDoctor(doctorId string) (string, error) {
	resp, err := a.doctorRepo.DeleteDoctor(doctorId) // Assuming Delete is implemented in PatientRepository
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// AddPatient - Admin can add a patient
func (a *adminService) AddPatient(patient domain.Patient) (string, error) {
	resp, err := a.patientRepo.SignUp(patient)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// DeletePatient - Admin can delete a patient
func (a *adminService) DeletePatient(patientID string) (string, error) {
	resp, err := a.patientRepo.DeletePatient(patientID) // Assuming Delete is implemented in PatientRepository
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// BlockPatient - Admin can block a patient
func (a *adminService) BlockPatient(patientID string, reason string) (string, error) {
	resp, err := a.patientRepo.Block(patientID, reason) // Assuming Block is implemented in PatientRepository
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (s *adminService) ListDoctors() ([]domain.Doctor, error) {
    return s.doctorRepo.ListDoctors()
}

// ListPatients returns a list of all patients
func (s *adminService) ListPatients() ([]domain.Patient, error) {
    return s.patientRepo.ListPatients()
}