package service

import (
	"github.com/go-redis/redis"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/repository"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/utils"
	"github.com/sirupsen/logrus"
)

type PatientService interface {
	SignIn(patient domain.Patient) (string, error)
	SignUp(patient domain.Patient) (string, error)
	SignUpVerify(token string) (string, error)
	GetProfile(patientId string) (domain.Patient, error)
	UpdateProfile(patient domain.Patient) error
	AddPrescription(data domain.PatientPrescription) error
	GetPrescription(patientId, query string) ([]domain.PatientPrescription, error)
	GetPatientCount() (int, error)
}

type patientService struct {
	repo   repository.PatientRepository
	logger *logrus.Logger
}

func NewPatientService(repo repository.PatientRepository, logger *logrus.Logger) PatientService {
	return &patientService{repo: repo, logger: logger}
}

// SignIn logic for patient
func (p *patientService) SignIn(patient domain.Patient) (string, error) {
	p.logger.WithFields(logrus.Fields{
		"function": "PatientSignIn",
		"id":       utils.MaskEmail(patient.Email),
	}).Info("Attempting patient login")

	resp, err := p.repo.SignIn(patient)
	if err != nil {
		p.logger.WithFields(logrus.Fields{
			"function": "PatientSignIn",
			"error":    err.Error(),
			"id":       utils.MaskEmail(patient.Email),
		}).Error("Login failed for patient")
		return resp, err
	}

	p.logger.WithFields(logrus.Fields{
		"function": "PatientSignIn",
		"status":   "success",
		"id":       utils.MaskEmail(patient.Email),
	}).Info("Patient login successful")
	return resp, nil
}
func (p *patientService) SignUp(patient domain.Patient) (string, error) {
	p.logger.WithFields(logrus.Fields{
		"function": "PatientSignUp",
		"email":    utils.MaskEmail(patient.Email),
	}).Info("Attempting patient signup")

	resp, err := p.repo.SignUp(patient)
	if err != nil {
		p.logger.WithFields(logrus.Fields{
			"function": "PatientSignUp",
			"error":    err.Error(),
			"email":    utils.MaskEmail(patient.Email),
		}).Error("Signup failed for patient")
		return resp, err
	}

	// Verification email sent
	resp, err = utils.SignUpverify(patient.Email)
	if err != nil {
		p.logger.WithFields(logrus.Fields{
			"function": "PatientSignUp",
			"error":    err.Error(),
			"email":    utils.MaskEmail(patient.Email),
		}).Error("Verification email failed to send")
		return "Failed to send verification email", err
	}

	p.logger.WithFields(logrus.Fields{
		"function": "PatientSignUp",
		"status":   "success",
		"email":    utils.MaskEmail(patient.Email),
	}).Info("Patient signup and verification email sent successfully")
	return resp, nil
}

func (p *patientService) UpdateProfile(patient domain.Patient) error {
	err := p.repo.UpdateProfile(patient)
	if err != nil {
		p.logger.WithFields(logrus.Fields{
			"function": "UpdateProfile",
			"error":    err.Error(),
			"id":       patient.ID,
		}).Error("Failed to update patient profile")
	}
	return err
}

func (p *patientService) GetProfile(patientId string) (domain.Patient, error) {
	patient, err := p.repo.GetProfile(patientId)
	if err != nil {
		p.logger.WithFields(logrus.Fields{
			"function":  "GetProfile",
			"error":     err.Error(),
			"patientId": patientId,
		}).Error("Failed to retrieve patient profile")
	}
	return patient, err
}

// GetPatientCount retrieves total number of patients (log only if useful for analytics/debugging)
func (p *patientService) GetPatientCount() (int, error) {
	count, err := p.repo.GetPatientCount()
	if err != nil {
		p.logger.WithFields(logrus.Fields{
			"function": "GetPatientCount",
			"error":    err.Error(),
		}).Error("Failed to retrieve patient count")
	}
	return count, err
}
func (p *patientService) SignUpVerify(token string) (string, error) {
	email, err := utils.Rdb.Get(token).Result()
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

func (p *patientService) AddPrescription(data domain.PatientPrescription) error {
	return p.repo.StorePatientPrescription(data)
}
func (p *patientService) GetPrescription(patientId, query string) ([]domain.PatientPrescription, error) {
	return p.repo.GetPrescriptions(patientId, query)
}
