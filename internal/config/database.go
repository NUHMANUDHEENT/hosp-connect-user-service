package config

import (
	"log"
	"os"

	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect with postgres......")
	}
	err = db.AutoMigrate(&domain.Admin{}, &domain.Doctor{}, &domain.Patient{},&domain.DoctorTokens{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
