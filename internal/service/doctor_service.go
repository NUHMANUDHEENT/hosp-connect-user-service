package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/repository"
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
}
type doctorService struct {
	repo repository.DoctorRepository
}

func NewDoctorService(repo repository.DoctorRepository) DoctorService {
	return &doctorService{repo: repo}
}
func (d *doctorService) DoctorSignin(email, password string) (string, error) {
	resp, err := d.repo.SignInValidate(email, password)
	fmt.Println("doctor id", resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (d *doctorService) GetProfile(email string) (domain.Doctor, error) {
	return d.repo.GetProfile(email)
}

func (d *doctorService) UpdateProfile(doctor domain.Doctor) error {
	return d.repo.UpdateProfile(doctor)
}

func (d *doctorService) StoreAccessToken(ctx context.Context, email string, token *oauth2.Token) error {
	return d.repo.StoreAccessToken(ctx, email, token)
}

func (d *doctorService) GetAccessToken(ctx context.Context, doctorID string) (*oauth2.Token, error) {
	return d.repo.GetAccessToken(ctx, doctorID)
}
func (d *doctorService) FetchAndStoreDoctorAvailability(ctx context.Context, doctorID string) ([]domain.AvailabilitySlot, error) {
	token, err := d.repo.GetAccessToken(ctx, doctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get doctor token: %w", err)
	}

	// Create a new OAuth2 client using the stored token
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	// Fetch Google Calendar service
	svc, err := calendar.New(oauthClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	// Define a list to hold the weekly availability
	var availability []domain.AvailabilitySlot

	// Define the time range (next 7 days)
	now := time.Now()
	oneWeekLater := now.AddDate(0, 0, 7).Format(time.RFC3339)

	// Fetch events from the primary calendar (adjust calendar ID as necessary)
	events, err := svc.Events.List("primary"). // Replace with actual calendar ID if needed
							TimeMin(now.Format(time.RFC3339)).
							TimeMax(oneWeekLater).
							SingleEvents(true).
							OrderBy("startTime").
							Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}

	// Process events and build availability
	for _, event := range events.Items {
		var startTime, endTime time.Time
		var err error

		if event.Start.DateTime != "" {
			startTime, err = time.Parse(time.RFC3339, event.Start.DateTime)
			if err != nil {
				return nil, fmt.Errorf("failed to parse event start time: %w", err)
			}
		} else {
			startTime, err = time.Parse("2006-01-02", event.Start.Date)
			if err != nil {
				return nil, fmt.Errorf("failed to parse event start date: %w", err)
			}
		}

		if event.End.DateTime != "" {
			endTime, err = time.Parse(time.RFC3339, event.End.DateTime)
			if err != nil {
				return nil, fmt.Errorf("failed to parse event end time: %w", err)
			}
		} else {
			endTime, err = time.Parse("2006-01-02", event.End.Date)
			if err != nil {
				return nil, fmt.Errorf("failed to parse event end date: %w", err)
			}
		}
		eventType := event.EventType

		availability = append(availability, domain.AvailabilitySlot{
			DoctorID:  doctorID,
			EventType: eventType,
			StartTime: startTime,
			EndTime:   endTime,
		})
		// }
	}

	// Store the availability in the database
	if err := d.repo.StoreDoctorSchedules(availability); err != nil {
		return nil, fmt.Errorf("failed to store availability: %w", err)
	}

	return availability, nil
}
