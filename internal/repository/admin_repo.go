package repository

import (
	"fmt"
	"os"

	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminRepository interface {
	SignIn(admin domain.Admin) (string, error)
	SignUp(admin domain.Admin) (string, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{
		db: db,
	}
}
func (a *adminRepository) SignIn(admin domain.Admin) (string, error) {
	var adminCheck domain.Admin
	if err := a.db.First(&adminCheck, "email = ?", admin.Email).Error; err != nil {
		if os.Getenv("ADMIN_EMAIL") == admin.Email && os.Getenv("ADMIN_PASSWORD") == admin.Password {
			password, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 12)
			if err != nil {
				return "Failed to hash password", err
			}
			admin.Role = "admin"
			admin.Password = string(password)
			if err := a.db.Create(&admin).Error; err != nil {
				return "Email already exists", err
			}
			return "Admin Sign in and new admin created successfully", nil
		}
		return "Email or password is incorrect", err
	}
	// Hashed password comparison
	if err := bcrypt.CompareHashAndPassword([]byte(adminCheck.Password), []byte(admin.Password)); err != nil {
		return "Email or password is incorrect", err
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(adminCheck.Password), []byte(admin.Password)); err != nil {
	// 	return "Email or password is incorrect", err
	// }
	fmt.Println("no error")
	return "Successfully login", nil
}
func (a *adminRepository) SignUp(admin domain.Admin) (string, error) {
	admin.Role = "admin"
	// Hash the password
	password, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 12)
	if err != nil {
		return "Failed to hash password", err
	}
	admin.Password = string(password)
	if err := a.db.Create(&admin).Error; err != nil {
		return "Email already exists", err
	}
	return "admin registered successfully", nil
}
