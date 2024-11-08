package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/repository"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

type DoctorService interface {
	DoctorSignin(email, password string) (string, error)
	GetProfile(email string) (domain.Doctor, error)
	UpdateProfile(doctor domain.Doctor) error
	StoreAccessToken(ctx context.Context, email string, token *oauth2.Token) error
	GetAccessToken(ctx context.Context, doctorID string) (*oauth2.Token, error)
	FetchAndStoreDoctorAvailability(ctx context.Context, doctorID string) ([]domain.AvailabilitySlot, error)
	GetAvailability(categoryId int32, reqDateTime time.Time) ([]domain.AvailabilitySlot, error)
	GetAvailabilityByDoctorId(doctorId string) ([]domain.AvailableDates, error)
	GetDoctorCount() (int, error)
}
type doctorService struct {
	repo   repository.DoctorRepository
	logger *logrus.Logger
}

func NewDoctorService(repo repository.DoctorRepository, logger *logrus.Logger) DoctorService {
	return &doctorService{repo: repo, logger: logger}
}

// DoctorSignin function with logging
func (d *doctorService) DoctorSignin(email, password string) (string, error) {
	d.logger.WithFields(logrus.Fields{
		"function": "DoctorSignin",
		"email":    utils.MaskEmail(email),
	}).Info("Attempting doctor sign-in")

	resp, err := d.repo.SignInValidate(email, password)
	if err != nil {
		d.logger.WithFields(logrus.Fields{
			"function": "DoctorSignin",
			"error":    err.Error(),
			"email":    utils.MaskEmail(email),
		}).Error("Failed to sign in doctor")
		return "", err
	}

	d.logger.WithFields(logrus.Fields{
		"function": "DoctorSignin",
		"doctorID": resp,
		"status":   "success",
	}).Info("Doctor signed in successfully")
	return resp, nil
}

// GetProfile function with logging
func (d *doctorService) GetProfile(email string) (domain.Doctor, error) {
	d.logger.WithFields(logrus.Fields{
		"function": "GetProfile",
		"id":       email,
	}).Info("Fetching doctor profile")

	profile, err := d.repo.GetProfile(email)
	if err != nil {
		d.logger.WithFields(logrus.Fields{
			"function": "GetProfile",
			"error":    err.Error(),
			// "email":    config.MaskEmail(email),
		}).Error("Failed to fetch doctor profile")
		return domain.Doctor{}, err
	}

	d.logger.WithFields(logrus.Fields{
		"function": "GetProfile",
		"doctorID": profile.ID,
		"status":   "success",
	}).Info("Fetched doctor profile successfully")
	return profile, nil
}

// UpdateProfile function with logging
func (d *doctorService) UpdateProfile(doctor domain.Doctor) error {
	d.logger.WithFields(logrus.Fields{
		"function": "UpdateProfile",
		"doctorID": doctor.ID,
	}).Info("Attempting to update doctor profile")
	log.Println(doctor)
	err := d.repo.UpdateProfile(doctor)
	if err != nil {
		d.logger.WithFields(logrus.Fields{
			"function": "UpdateProfile",
			"error":    err.Error(),
			"doctorID": doctor.ID,
		}).Error("Failed to update doctor profile")
		return err
	}

	d.logger.WithFields(logrus.Fields{
		"function": "UpdateProfile",
		"doctorID": doctor.ID,
		"status":   "success",
	}).Info("Doctor profile updated successfully")
	return nil
}

// StoreAccessToken function with logging
func (d *doctorService) StoreAccessToken(ctx context.Context, email string, token *oauth2.Token) error {
	d.logger.WithFields(logrus.Fields{
		"function": "StoreAccessToken",
		"email":    utils.MaskEmail(email),
	}).Info("Storing access token for doctor")

	err := d.repo.StoreAccessToken(ctx, email, token)
	if err != nil {
		d.logger.WithFields(logrus.Fields{
			"function": "StoreAccessToken",
			"error":    err.Error(),
			"email":    utils.MaskEmail(email),
		}).Error("Failed to store access token")
		return err
	}

	d.logger.WithFields(logrus.Fields{
		"function": "StoreAccessToken",
		"email":    utils.MaskEmail(email),
		"status":   "success",
	}).Info("Access token stored successfully")
	return nil
}

// GetAccessToken function with logging
func (d *doctorService) GetAccessToken(ctx context.Context, doctorID string) (*oauth2.Token, error) {
	d.logger.WithFields(logrus.Fields{
		"function": "GetAccessToken",
		"doctorID": doctorID,
	}).Info("Fetching access token for doctor")

	token, err := d.repo.GetAccessToken(ctx, doctorID)
	if err != nil {
		d.logger.WithFields(logrus.Fields{
			"function": "GetAccessToken",
			"error":    err.Error(),
			"doctorID": doctorID,
		}).Error("Failed to fetch access token")
		return nil, err
	}

	d.logger.WithFields(logrus.Fields{
		"function": "GetAccessToken",
		"doctorID": doctorID,
		"status":   "success",
	}).Info("Access token fetched successfully")
	return token, nil
}

// FetchAndStoreDoctorAvailability function with logging
func (d *doctorService) FetchAndStoreDoctorAvailability(ctx context.Context, doctorID string) ([]domain.AvailabilitySlot, error) {
	d.logger.WithFields(logrus.Fields{
		"function": "FetchAndStoreDoctorAvailability",
		"doctorID": doctorID,
	}).Info("Fetching doctor availability")

	token, err := d.repo.GetAccessToken(ctx, doctorID)
	if err != nil {
		d.logger.WithFields(logrus.Fields{
			"function": "FetchAndStoreDoctorAvailability",
			"error":    err.Error(),
			"doctorID": doctorID,
		}).Error("Failed to get doctor token")
		return nil, err
	}

	// Create a new OAuth2 client using the stored token
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	// Fetch Google Calendar service
	svc, err := calendar.New(oauthClient)
	if err != nil {
		d.logger.WithFields(logrus.Fields{
			"function": "FetchAndStoreDoctorAvailability",
			"error":    err.Error(),
		}).Error("Failed to create calendar service")
		return nil, err
	}

	// Define a list to hold the weekly availability
	var availability []domain.AvailabilitySlot

	// Define the time range (next 7 days)
	now := time.Now()
	oneWeekLater := now.AddDate(0, 0, 7).Format(time.RFC3339)

	// Fetch events from the primary calendar
	events, err := svc.Events.List("primary").
		TimeMin(now.Format(time.RFC3339)).
		TimeMax(oneWeekLater).
		SingleEvents(true).
		OrderBy("startTime").
		Do()
	if err != nil {
		d.logger.WithFields(logrus.Fields{
			"function": "FetchAndStoreDoctorAvailability",
			"error":    err.Error(),
		}).Error("Failed to fetch events from calendar")
		return nil, err
	}

	// Process events and build availability
	for _, event := range events.Items {
		var startTime, endTime time.Time
		var err error

		if event.Start.DateTime != "" {
			startTime, err = time.Parse(time.RFC3339, event.Start.DateTime)
			if err != nil {
				d.logger.WithFields(logrus.Fields{
					"function": "FetchAndStoreDoctorAvailability",
					"error":    err.Error(),
				}).Error("Failed to parse event start time")
				return nil, err
			}
		} else {
			startTime, err = time.Parse("2006-01-02", event.Start.Date)
			if err != nil {
				d.logger.WithFields(logrus.Fields{
					"function": "FetchAndStoreDoctorAvailability",
					"error":    err.Error(),
				}).Error("Failed to parse event start date")
				return nil, err
			}
		}

		if event.End.DateTime != "" {
			endTime, err = time.Parse(time.RFC3339, event.End.DateTime)
			if err != nil {
				d.logger.WithFields(logrus.Fields{
					"function": "FetchAndStoreDoctorAvailability",
					"error":    err.Error(),
				}).Error("Failed to parse event end time")
				return nil, err
			}
		} else {
			endTime, err = time.Parse("2006-01-02", event.End.Date)
			if err != nil {
				d.logger.WithFields(logrus.Fields{
					"function": "FetchAndStoreDoctorAvailability",
					"error":    err.Error(),
				}).Error("Failed to parse event end date")
				return nil, err
			}
		}
		eventType := event.EventType

		availability = append(availability, domain.AvailabilitySlot{
			DoctorID:  doctorID,
			EventType: eventType,
			StartTime: startTime,
			EndTime:   endTime,
		})
	}

	// Store the availability in the database
	if err := d.repo.StoreDoctorSchedules(availability); err != nil {
		d.logger.WithFields(logrus.Fields{
			"function": "FetchAndStoreDoctorAvailability",
			"error":    err.Error(),
		}).Error("Failed to store availability")
		return nil, err
	}

	d.logger.WithFields(logrus.Fields{
		"function": "FetchAndStoreDoctorAvailability",
		"doctorID": doctorID,
		"status":   "success",
	}).Info("Doctor availability stored successfully")
	return availability, nil
}
func (d *doctorService) GetAvailability(categoryId int32, reqDateTime time.Time) ([]domain.AvailabilitySlot, error) {
	resp, err := d.repo.GetAvailabilityByCategory(categoryId, reqDateTime)
	if err != nil {
		return []domain.AvailabilitySlot{}, err
	}
	return resp, nil

}
func (d *doctorService) GetAvailabilityByDoctorId(doctorId string) ([]domain.AvailableDates, error) {
	resp, err := d.repo.GetAvailabilityByDoctorId(doctorId)
	fmt.Println("availabiluty=========", resp)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
func (p *doctorService) GetDoctorCount() (int, error) {
	return p.repo.GetDoctorCount()
}
