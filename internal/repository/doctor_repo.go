package repository

import (
	"context"

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
	StoreAccessToken(ctx context.Context, doctorID string, token *oauth2.Token) error
	GetAccessToken(ctx context.Context, doctorID string) (*oauth2.Token, error)
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

	return "Successfully logged in", nil
}
func (d *doctorRepository) SignUpStore(doctordetails domain.Doctor) (string, error) {
	doctordetails.Role = "doctor"
	password, err := bcrypt.GenerateFromPassword([]byte(doctordetails.Password), 12)
	if err != nil {
		return "Failed to hash password", err
	}
	doctordetails.Password = string(password)
	if err := d.db.Create(&doctordetails).Error; err != nil {
		return "Failed to create doctor", err
	}

	return "Doctor Registed successfully", nil

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
	if err := d.db.Where("email = ?", email).First(&doctor).Error; err != nil {
		return doctor, err
	}
	return doctor, nil
}

// UpdateProfile updates a doctor's profile
func (d *doctorRepository) UpdateProfile(doctor domain.Doctor) error {
	if err := d.db.Where("email=?", doctor.Email).Model(&doctor).Updates(doctor).Error; err != nil {
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

func (r *doctorRepository) StoreAccessToken(ctx context.Context, doctorID string, token *oauth2.Token) error {
	tokenstore := domain.DoctorTokens{
		DoctorId:     doctorID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}
	if err := r.db.Create(&tokenstore).Error; err != nil {
		return err
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
