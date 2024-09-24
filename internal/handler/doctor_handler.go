package handler

import (
	"context"
	"time"

	pb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/doctor"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/service"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/anypb"
)

type DoctorServiceClient struct {
	pb.UnimplementedDoctorServiceServer
	service service.DoctorService
}

func NewDoctorHandler(service service.DoctorService) *DoctorServiceClient {
	return &DoctorServiceClient{
		service: service,
	}
}

func (d *DoctorServiceClient) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.StandardResponse, error) {
	resp, err := d.service.DoctorSignin(req.Email, req.Password)
	if err != nil {
		return &pb.StandardResponse{Error: resp, StatusCode: 401, Status: "fail"}, nil
	}
	return &pb.StandardResponse{Message: resp, StatusCode: 200, Status: "success"}, nil
}

func (d *DoctorServiceClient) StoreAccessToken(ctx context.Context, req *pb.StoreAccessTokenRequest) (*pb.StandardResponse, error) {
	token := &oauth2.Token{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
		Expiry:       time.Now(), // Assuming expiry is passed as string, you can convert it to `time.Time`
	}

	err := d.service.StoreAccessToken(ctx, req.DoctorId, token)
	if err != nil {
		return &pb.StandardResponse{
			Status:     "fail",
			StatusCode: 500,
			Error:      err.Error(),
		}, nil
	}

	return &pb.StandardResponse{
		Status:     "success",
		StatusCode: 200,
		Message:    "Token stored successfully",
	}, nil
}

func (d *DoctorServiceClient) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	doctor, err := d.service.GetProfile(req.Email)
	if err != nil {
		return &pb.GetProfileResponse{
			Status:     "fail",
			StatusCode: 404,
			Error:      err.Error(),
		}, nil
	}

	return &pb.GetProfileResponse{
		Email:            doctor.Email,
		Name:             doctor.Name,
		Phone:            int32(doctor.Phone),
		SpecializationId: int32(doctor.SpecilazationId),
		Status:           "success",
		StatusCode:       200,
	}, nil
}

func (d *DoctorServiceClient) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.StandardResponse, error) {
	err := d.service.UpdateProfile(domain.Doctor{
		DoctorId:        req.Doctor.DoctorId,
		Name:            req.Doctor.Name,
		Email:           req.Doctor.Email,
		SpecilazationId: int(req.Doctor.SpecializationId),
	})
	if err != nil {
		return &pb.StandardResponse{
			Status:     "fail",
			StatusCode: 400,
		}, err
	}
	return &pb.StandardResponse{
		Message:    "Profile updated successfully",
		Status:     "success",
		StatusCode: 200,
	}, nil
}

// GetAccessToken handles fetching OAuth token for a doctor
func (d *DoctorServiceClient) GetAccessToken(ctx context.Context, req *pb.GetAccessTokenRequest) (*pb.StandardResponse, error) {
	token, err := d.service.GetAccessToken(ctx, req.DoctorId)
	if err != nil {
		return nil, err
	}
	tokenStruct := &pb.StoreAccessTokenRequest{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	anypbstruct, err := anypb.New(tokenStruct)
	if err != nil {
		return nil, err
	}
	return &pb.StandardResponse{
		Status:     "success",
		StatusCode: 200,
		Message:    "Token stored successfully",
		Data:       anypbstruct,
	}, nil
}
