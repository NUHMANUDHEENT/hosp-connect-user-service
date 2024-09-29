package domain

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	AdminId  string
	Name     string
	Email    string `gorm:"uniqueIndex"`
	Password string
	Role     string
}
