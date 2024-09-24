package domain

import "gorm.io/gorm"

type Admin struct {
    gorm.Model
    Name  string 
    Email string `gorm:"uniqueIndex"`
    Password string
    Role string
}