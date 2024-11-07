package service

import (
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/config"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/repository"
	"github.com/sirupsen/logrus"
)

type AdminService interface {
	SignIn(admin domain.Admin) (string, error)
	SignUp(email, name, password string) (string, error)
	AddDoctor(email, name, password string, phone, specializationId int32) (string, error)
	DeleteDoctor(doctorId string) (string, error)
	AddPatient(patient domain.Patient) (string, error)
	DeletePatient(patientID string) (string, error)
	BlockPatient(patientID string, reason string) (string, error)
	ListDoctors() ([]domain.Doctor, error)
	ListPatients() ([]domain.Patient, error)
}
type adminService struct {
	repo        repository.AdminRepository
	doctorRepo  repository.DoctorRepository
	patientRepo repository.PatientRepository
	logger      *logrus.Logger
}

// NewAdminService constructor
func NewAdminService(repo repository.AdminRepository, docrepo repository.DoctorRepository, patientRepo repository.PatientRepository, logger *logrus.Logger) AdminService {
	return &adminService{repo: repo,
		doctorRepo:  docrepo,
		patientRepo: patientRepo,
		logger:      logger,
	}
}

// SignIn logic for admin
func (a *adminService) SignIn(admin domain.Admin) (string, error) {
	// Log SignIn attempt
	a.logger.WithFields(logrus.Fields{
		"function": "AdminSignIn",
		"email":    config.MaskEmail(admin.Email),
	}).Info("Attempting admin login")

	resp, err := a.repo.SignIn(admin)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"function": "AdminSignIn",
			"error":    err.Error(),
			"email":    config.MaskEmail(admin.Email),
		}).Error("Admin login failed")
		return resp, err
	}

	a.logger.WithFields(logrus.Fields{
		"function": "AdminSignIn",
		"status":   "success",
		"email":    config.MaskEmail(admin.Email),
	}).Info("Admin login successful")
	return resp, nil
}

func (a *adminService) SignUp(email, name, password string) (string, error) {
	a.logger.WithFields(logrus.Fields{
		"function": "AdminSignUp",
		"email":    config.MaskEmail(email),
	}).Info("Attempting to create a new admin account")

	resp, err := a.repo.SignUp(domain.Admin{Email: email, Name: name, Password: password})
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"function": "AdminSignUp",
			"error":    err.Error(),
			"email":    config.MaskEmail(email),
		}).Error("Failed to create admin account")
		return resp, err
	}

	a.logger.WithFields(logrus.Fields{
		"function": "AdminSignUp",
		"status":   "success",
		"email":    config.MaskEmail(email),
	}).Info("Admin account created successfully")
	return resp, nil
}

func (a *adminService) AddDoctor(email, name, password string, phone, specializationId int32) (string, error) {
	a.logger.WithFields(logrus.Fields{
		"function":         "AddDoctor",
		"email":            config.MaskEmail(email),
		"phone":            phone,
		"specializationId": specializationId,
	}).Info("Attempting to add a new doctor")

	resp, err := a.doctorRepo.SignUpStore(domain.Doctor{
		Email:            email,
		Name:             name,
		Password:         password,
		SpecializationId: int(specializationId),
		Phone:            int(phone),
	})
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"function": "AddDoctor",
			"error":    err.Error(),
			"email":    config.MaskEmail(email),
		}).Error("Failed to add doctor")
		return resp, err
	}

	a.logger.WithFields(logrus.Fields{
		"function": "AddDoctor",
		"status":   "success",
		"email":    config.MaskEmail(email),
	}).Info("Doctor added successfully")
	return resp, nil
}

func (a *adminService) DeleteDoctor(doctorId string) (string, error) {
	a.logger.WithFields(logrus.Fields{
		"function": "DeleteDoctor",
		"doctorId": doctorId,
	}).Info("Attempting to delete doctor")

	resp, err := a.doctorRepo.DeleteDoctor(doctorId)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"function": "DeleteDoctor",
			"error":    err.Error(),
			"doctorId": doctorId,
		}).Error("Failed to delete doctor")
		return resp, err
	}

	a.logger.WithFields(logrus.Fields{
		"function": "DeleteDoctor",
		"status":   "success",
		"doctorId": doctorId,
	}).Info("Doctor deleted successfully")
	return resp, nil
}

func (a *adminService) AddPatient(patient domain.Patient) (string, error) {
	a.logger.WithFields(logrus.Fields{
		"function": "AddPatient",
		"email":    config.MaskEmail(patient.Email),
		"name":     patient.Name,
	}).Info("Attempting to add a new patient")

	resp, err := a.patientRepo.SignUp(patient)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"function": "AddPatient",
			"error":    err.Error(),
			"email":    config.MaskEmail(patient.Email),
		}).Error("Failed to add patient")
		return resp, err
	}

	a.logger.WithFields(logrus.Fields{
		"function": "AddPatient",
		"status":   "success",
		"email":    config.MaskEmail(patient.Email),
	}).Info("Patient added successfully")
	return resp, nil
}

func (a *adminService) DeletePatient(patientID string) (string, error) {
	a.logger.WithFields(logrus.Fields{
		"function":  "DeletePatient",
		"patientId": patientID,
	}).Info("Attempting to delete patient")

	resp, err := a.patientRepo.DeletePatient(patientID)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"function":  "DeletePatient",
			"error":     err.Error(),
			"patientId": patientID,
		}).Error("Failed to delete patient")
		return resp, err
	}

	a.logger.WithFields(logrus.Fields{
		"function":  "DeletePatient",
		"status":    "success",
		"patientId": patientID,
	}).Info("Patient deleted successfully")
	return resp, nil
}

func (a *adminService) BlockPatient(patientID string, reason string) (string, error) {
	a.logger.WithFields(logrus.Fields{
		"function":  "BlockPatient",
		"patientId": patientID,
		"reason":    reason,
	}).Info("Attempting to block patient")

	resp, err := a.patientRepo.Block(patientID, reason)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"function":  "BlockPatient",
			"error":     err.Error(),
			"patientId": patientID,
		}).Error("Failed to block patient")
		return resp, err
	}

	a.logger.WithFields(logrus.Fields{
		"function":  "BlockPatient",
		"status":    "success",
		"patientId": patientID,
	}).Info("Patient blocked successfully")
	return resp, nil
}

func (s *adminService) ListDoctors() ([]domain.Doctor, error) {
	s.logger.WithFields(logrus.Fields{
		"function": "ListDoctors",
	}).Info("Fetching list of doctors")

	doctors, err := s.doctorRepo.ListDoctors()
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"function": "ListDoctors",
			"error":    err.Error(),
		}).Error("Failed to retrieve list of doctors")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"function": "ListDoctors",
		"status":   "success",
		"count":    len(doctors),
	}).Info("List of doctors retrieved successfully")
	return doctors, nil
}

func (s *adminService) ListPatients() ([]domain.Patient, error) {
	s.logger.WithFields(logrus.Fields{
		"function": "ListPatients",
	}).Info("Fetching list of patients")

	patients, err := s.patientRepo.ListPatients()
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"function": "ListPatients",
			"error":    err.Error(),
		}).Error("Failed to retrieve list of patients")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"function": "ListPatients",
		"status":   "success",
		"count":    len(patients),
	}).Info("List of patients retrieved successfully")
	return patients, nil
}
