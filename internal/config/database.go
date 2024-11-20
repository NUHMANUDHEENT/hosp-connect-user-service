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
	log.Println("Database url :", dsn)
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := gorm.Open(postgres.Open("postgres://postgres:Nuhman%40456@postgres-db-service:5432/user_svc"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect with postgres......",err)
	}
	err = db.AutoMigrate(&domain.Admin{}, &domain.Doctor{}, &domain.Patient{}, &domain.DoctorTokens{}, &domain.AvailabilitySlot{}, &domain.PatientPrescription{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
