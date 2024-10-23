package repository

import (
	"fmt"

	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
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
		return "Email or password is incorrect", err
	}
	if admin.Password != adminCheck.Password {
		return "Email or password is incorrect", nil
	}
	// if err := bcrypt.CompareHashAndPassword([]byte(adminCheck.Password), []byte(admin.Password)); err != nil {
	// 	return "Email or password is incorrect", err
	// }
	fmt.Println("no error")
	return "Successfully login", nil
}
func (a *adminRepository) SignUp(admin domain.Admin) (string, error){
	admin.Role = "admin"
	if err := a.db.Create(&admin).Error; err != nil {
		return "Email already exists", err
	}
	return "admin registered successfully", nil
}