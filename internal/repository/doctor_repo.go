package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	StoreDoctorSchedules(schedules []domain.AvailabilitySlot,doctorid string) error
	GetAvailabilityByCategory(categoryId int32, reqDateTime time.Time) ([]domain.AvailabilitySlot, error)
	GetAvailabilityByDoctorId(doctorId string) ([]domain.AvailableDates, error)
	GetDoctorCount() (int, error)
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

	if err := d.db.Create(&doctordetails).Error; err != nil {
		return "Email already exists", err
	}

	return "Doctor registered successfully", nil
}
func (p *doctorRepository) DeleteDoctor(doctorId string) (string, error) {
	var doctor domain.Doctor
	if err := p.db.First(&doctor,"doctor_id", doctorId).Error; err != nil {
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
	if err := d.db.Where("doctor_id =?", doctor.DoctorId).Model(&doctor).Updates(doctor).Error; err != nil {
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
		DoctorId:     doctorDetails.DoctorId,
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
func (r *doctorRepository) StoreDoctorSchedules(schedules []domain.AvailabilitySlot,doctorid string) error {
	doctors := domain.Doctor{}
	if err := r.db.Where("doctor_id = ?", doctorid).First(&doctors).Error; err != nil {
		return err
	}
	for i, _ := range schedules {
		schedules[i].DoctorName = doctors.Name
	}
	result := r.db.Create(&schedules)

	if result.Error != nil {
		return fmt.Errorf("failed to store schedules: %w", result.Error)
	}

	return nil
}


// doctorRepository.go
func (r *doctorRepository) GetAvailabilityByCategory(categoryId int32, reqDateTime time.Time) ([]domain.AvailabilitySlot, error) {
	var availabilities []domain.AvailabilitySlot

	// Query doctors with the provided CategoryId
	doctors := []domain.Doctor{}
	if err := r.db.Where("specialization_id = ?", categoryId).Find(&doctors).Error; err != nil {
		return nil, err
	}
	// Extract Doctor IDs from the result
	doctorIds := make([]string, len(doctors))
	for i, doctor := range doctors {
		doctorIds[i] = doctor.DoctorId
	}
	fmt.Println("Requested doctors", doctorIds, "for date:", reqDateTime)

	// Fetch unavailable slots for doctors in the category (those who have scheduled unavailability during the requested time)
	unavailableDoctors := []domain.AvailabilitySlot{}
	if err := r.db.Where("doctor_id IN (?)", doctorIds).
		Where("start_time <= ?", reqDateTime).
		Where("end_time >= ?", reqDateTime).
		Find(&unavailableDoctors).Error; err != nil {
		return nil, err
	}

	// Create a set of unavailable doctor IDs for quick lookup
	unavailableDoctorIds := map[string]bool{}
	for _, slot := range unavailableDoctors {
		unavailableDoctorIds[slot.DoctorID] = true
	}

	// Filter out unavailable doctors
	for _, doctor := range doctors {
		if !unavailableDoctorIds[doctor.DoctorId] {
			availabilities = append(availabilities, domain.AvailabilitySlot{
				DoctorID:   doctor.DoctorId,
				DoctorName: doctor.Name,
			})
		}
	}

	if len(availabilities) == 0 {
		fmt.Println("No doctors are available at the requested time")
		return nil, nil
	}

	return availabilities, nil
}

func (r *doctorRepository) GetAvailabilityByDoctorId(doctorId string) ([]domain.AvailableDates, error) {
	var leaveDate []domain.AvailabilitySlot

	startime := time.Now()
	endtime := time.Now().AddDate(0, 0, 7)
	if err := r.db.Model(&domain.AvailabilitySlot{}).
		Where("doctor_id = ? AND start_time >? AND start_time<?", doctorId, startime, endtime).
		Find(&leaveDate).Error; err != nil {
		return nil, err
	}
	if leaveDate == nil {
		return nil, errors.New("no slot available on next 7 days")
	}
	available := Create7daysAvailability(leaveDate)

	return available, nil // Doctor is available
}
func Create7daysAvailability(leaveDates []domain.AvailabilitySlot) []domain.AvailableDates {
	var availability []domain.AvailableDates

	// Start from the current day and check the next 7 days
	startTime := time.Now()
	endTime := time.Now().AddDate(0, 0, 7)

	// Iterate through each day for the next 7 days
	for current := startTime; current.Before(endTime); current = current.AddDate(0, 0, 1) {
		// Skip Sundays
		if current.Weekday() == time.Sunday {
			continue
		}

		// Check if the current date is in the leave list
		isAvailable := "available"
		for _, leave := range leaveDates {
			if leave.StartTime.Year() == current.Year() &&
				leave.StartTime.YearDay() == current.YearDay() {
				isAvailable = "unavailable"
				break
			}
		}

		// Append the availability status for the current day
		availability = append(availability, domain.AvailableDates{
			DateTime:    current,
			IsAvailable: isAvailable,
		})
	}

	return availability
}
func (p *doctorRepository) GetDoctorCount() (int, error) {
	var count int64
	err := p.db.Model(&domain.Doctor{}).Count(&count).Error
	return int(count), err
}
