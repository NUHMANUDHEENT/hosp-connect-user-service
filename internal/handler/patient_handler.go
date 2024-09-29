package handler

import (
	"context"
	"log"

	pb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/patient"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/service"
)

type PatientServiceClient struct {
	pb.UnimplementedPatientServiceServer
	service service.PatientService
}

func NewPatientHandler(service service.PatientService) *PatientServiceClient {
	return &PatientServiceClient{
		service: service,
	}
}

// Patient SignIn Handler
func (p *PatientServiceClient) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	patientDetails := domain.Patient{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	log.Println("Patient signin request", patientDetails)

	resp, err := p.service.SignIn(patientDetails)
	if err != nil {
		return &pb.SignInResponse{
			Status:     "fail",
			Message:    "Invalid credentials, please try again.",
			StatusCode: 401,
		}, nil
	}
	return &pb.SignInResponse{
		PatientId:  resp,
		Status:     "success",
		Message:    "Successfully logged in",
		StatusCode: 200,
	}, nil
}

func (p *PatientServiceClient) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.StandardResponse, error) {
	patientDetails := domain.Patient{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Name:     req.GetName(),
		Phone:    int(req.Phone),
		Age:      req.GetAge(),
		Gender:   req.GetGender(),
		IsBlock:  false,
	}

	log.Println("Patient signup request", patientDetails)

	resp, err := p.service.SignUp(patientDetails)
	if err != nil {
		return &pb.StandardResponse{
			Status:     "fail",
			Error:      err.Error(),
			StatusCode: 400,
		}, nil
	}
	return &pb.StandardResponse{
		Status:     "success",
		Message:    resp,
		StatusCode: 200,
	}, nil
}
func (p *PatientServiceClient) SignUpVerify(ctx context.Context, req *pb.SignUpVerifyRequest) (*pb.StandardResponse, error) {
	resp, err := p.service.SignUpVerify(req.Token)
	if err != nil {
		return &pb.StandardResponse{
			Status:     "fail",
			Error:      resp,
			StatusCode: 400,
		}, nil
	}
	return &pb.StandardResponse{
		Status:     "success",
		Message:    resp,
		StatusCode: 200,
	}, nil
}
func (p *PatientServiceClient) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	patient, err := p.service.GetProfile(req.PatientId)
	if err != nil {
		return &pb.GetProfileResponse{
			Status:     "fail",
			StatusCode: 404,
			Error:      err.Error(),
		}, nil
	}

	return &pb.GetProfileResponse{
		Email:      patient.Email,
		Name:       patient.Name,
		Phone:      int32(patient.Phone),
		Age:        int32(patient.Age),
		Gender:     patient.Gender,
		Status:     "success",
		StatusCode: 200,
	}, nil
}

func (p *PatientServiceClient) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.StandardResponse, error) {
	err := p.service.UpdateProfile(domain.Patient{
		PatientID: req.Patient.PatientId,
		Name:      req.Patient.Name,
		Email:     req.Patient.Email,
		Phone:     int(req.Patient.Phone),
		Age:       req.Patient.Age,
		Gender:    req.Patient.Gender,
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
