package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type DoctorRepository interface {
	SignInValidate(email, password string) (string, error)
	SignUpStore(doctordetails domain.Doctor) (string, error)
	DeleteDoctor(doctorId string) (string, error)
	GetProfile(email string) (domain.Doctor, error)
	UpdateProfile(doctor domain.Doctor) error
	ListDoctors() ([]domain.Doctor, error)
	StoreAccessToken(ctx context.Context, email string, token *oauth2.Token) error
	GetAccessToken(ctx context.Context, doctorID string) (*oauth2.Token, error)
	StoreDoctorSchedules(schedules []domain.AvailabilitySlot) error
	CreateSpecialization(specialize domain.DoctorSpecialization) (string, error)
}
type doctorRepository struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepository{db: db}
}

func (d *doctorRepository) SignInValidate(email, password string) (string, error) {
	var doctorCheck domain.Doctor
	if err := d.db.First(&doctorCheck, "email = ?", email).Error; err != nil {
		return "Email or password is incorrect", err
	}

	// Hashed password comparison
	if err := bcrypt.CompareHashAndPassword([]byte(doctorCheck.Password), []byte(password)); err != nil {
		return "Email or password is incorrect", err
	}

	return doctorCheck.DoctorId, nil
}
func (d *doctorRepository) SignUpStore(doctordetails domain.Doctor) (string, error) {
	// Set the role as doctor
	doctordetails.Role = "doctor"

	// Hash the password
	password, err := bcrypt.GenerateFromPassword([]byte(doctordetails.Password), 12)
	if err != nil {
		return "Failed to hash password", err
	}
	doctordetails.Password = string(password)

	// Generate a unique DoctorId with prefix 'dc'
	uuid := uuid.New()
	doctordetails.DoctorId = fmt.Sprintf("dc-%s", uuid.String())

	// Store the doctor details in the database
	if err := d.db.Create(&doctordetails).Error; err != nil {
		return "Email already exists", err
	}

	return "Doctor registered successfully", nil
}
func (p *doctorRepository) DeleteDoctor(doctorId string) (string, error) {
	var doctor domain.Doctor
	if err := p.db.First(&doctor, doctorId).Error; err != nil {
		return "Doctor not found", err
	}
	if err := p.db.Delete(&doctor).Error; err != nil {
		return "Doctor deletion failed", err
	}
	return "Doctor deleted successfully", nil
}
func (d *doctorRepository) GetProfile(email string) (domain.Doctor, error) {
	var doctor domain.Doctor
	if err := d.db.Where("doctor_id = ?", email).First(&doctor).Error; err != nil {
		return doctor, err
	}
	return doctor, nil
}

// UpdateProfile updates a doctor's profile
func (d *doctorRepository) UpdateProfile(doctor domain.Doctor) error {
	if err := d.db.Where("doctor_id =?", doctor.Email).Model(&doctor).Updates(doctor).Error; err != nil {
		return err
	}
	return nil
}

// ListDoctors returns a list of all doctors (admin feature)
func (d *doctorRepository) ListDoctors() ([]domain.Doctor, error) {
	var doctors []domain.Doctor
	if err := d.db.Find(&doctors).Error; err != nil {
		return nil, err
	}
	return doctors, nil
}

func (r *doctorRepository) StoreAccessToken(ctx context.Context, email string, token *oauth2.Token) error {
	var doctorDetails domain.Doctor
	if err := r.db.Where("email = ?", email).First(&doctorDetails).Error; err != nil {
		return errors.New("doctor not found")
	}
	var tokenStore domain.DoctorTokens
	err := r.db.First(&tokenStore, "doctor_id=?", doctorDetails.DoctorId).Error
	tokenStore = domain.DoctorTokens{
		DoctorId:     tokenStore.DoctorId,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}
	if err == nil {
		if err := r.db.Where("doctor_id= ?", tokenStore.DoctorId).Model(&tokenStore).Updates(tokenStore).Error; err != nil {
			return errors.New("FAILED TO UPDATE")
		}
	} else {
		if err := r.db.Create(&tokenStore).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAccessToken retrieves the OAuth token from the database
func (r *doctorRepository) GetAccessToken(ctx context.Context, doctorID string) (*oauth2.Token, error) {
	var fetchToken domain.DoctorTokens
	if err := r.db.Where("doctor_id = ?", doctorID).First(&fetchToken).Error; err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  fetchToken.AccessToken,
		RefreshToken: fetchToken.RefreshToken,
		Expiry:       fetchToken.ExpiresAt,
	}

	return token, nil
}
func (r *doctorRepository) StoreDoctorSchedules(schedules []domain.AvailabilitySlot) error {
	result := r.db.Create(&schedules)

	if result.Error != nil {
		return fmt.Errorf("failed to store schedules: %w", result.Error)
	}

	return nil

}
func (r *doctorRepository) CreateSpecialization(specialize domain.DoctorSpecialization) (string, error) {
	if err := r.db.Create(&specialize).Error; err != nil {
		return "Category is already exist", err
	}
	return "Category created successfully", nil
}
