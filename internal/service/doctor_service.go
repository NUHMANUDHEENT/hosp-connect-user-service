package service

import (
	"context"
	"fmt"

	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/repository"
	"golang.org/x/oauth2"
)

type DoctorService interface {
	DoctorSignin(email, password string) (string, error)
	GetProfile(email string) (domain.Doctor, error)
	UpdateProfile(doctor domain.Doctor) error
	StoreAccessToken(ctx context.Context, doctorID string, token *oauth2.Token) error
	GetAccessToken(ctx context.Context, doctorID string) (*oauth2.Token, error)
}
type doctorService struct {
	repo repository.DoctorRepository
}

func NewDoctorService(repo repository.DoctorRepository) DoctorService {
	return &doctorService{repo: repo}
}
func (d *doctorService) DoctorSignin(email, password string) (string, error) {
	resp, err := d.repo.SignInValidate(email, password)
	fmt.Println(err)
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

func (d *doctorService) StoreAccessToken(ctx context.Context, doctorID string, token *oauth2.Token) error {
	return d.repo.StoreAccessToken(ctx, doctorID, token)
}

func (d *doctorService) GetAccessToken(ctx context.Context, doctorID string) (*oauth2.Token, error) {
	return d.repo.GetAccessToken(ctx, doctorID)
}
