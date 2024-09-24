package handler

import (
	"context"
	"log"

	pb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/admin"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-user-service/internal/service"
)

type AdminServiceClient struct {
	pb.UnimplementedAdminServiceServer
	service service.AdminService
}

func NewAdminHandler(service service.AdminService) *AdminServiceClient {
	return &AdminServiceClient{
		service: service,
	}
}

func (a *AdminServiceClient) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.StandardResponse, error) {
	admindetails := domain.Admin{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	log.Println("signin request", admindetails)
	resp, err := a.service.SignIn(admindetails)
	if err != nil {
		return &pb.StandardResponse{
			Status:     "fail",
			Error:      "Invalid credentials, please try again.",
			StatusCode: 401,
		}, nil

	}
	return &pb.StandardResponse{
		Status:     "success",
		Message:    resp,
		StatusCode: 200,
	}, nil
}
func (a *AdminServiceClient) AddDoctor(ctx context.Context, req *pb.AddDoctorRequest) (*pb.StandardResponse, error) {
	log.Println("Regiter doctor  from admin with ", req.Email, req.Name)
	resp, err := a.service.AddDoctor(req.Email, req.Name, req.Password, req.Phone, req.SpecializationId)
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
func (a *AdminServiceClient) DeleteDoctor(ctx context.Context, req *pb.DeleteDoctorRequest) (*pb.StandardResponse, error) {
	log.Println("Deleting patient with ID: ", req.DoctorId)
	resp, err := a.service.DeleteDoctor(req.GetDoctorId())
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


// Add Patient
func (a *AdminServiceClient) AddPatient(ctx context.Context, req *pb.AddPatientRequest) (*pb.StandardResponse, error) {
	patientDetails := domain.Patient{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Name:     req.GetName(),
		Phone:    int(req.GetPhone()),
	}
	log.Println("Adding patient: ", patientDetails)
	resp, err := a.service.AddPatient(patientDetails)
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

// Delete Patient
func (a *AdminServiceClient) DeletePatient(ctx context.Context, req *pb.DeletePatientRequest) (*pb.StandardResponse, error) {
	log.Println("Deleting patient with ID: ", req.PatientId)
	resp, err := a.service.DeletePatient(req.GetPatientId())
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

// Block Patient
func (a *AdminServiceClient) BlockPatient(ctx context.Context, req *pb.BlockPatientRequest) (*pb.StandardResponse, error) {
	log.Println("Blocking patient with ID: ", req.PatientId)
	resp, err := a.service.BlockPatient(req.GetPatientId(), req.GetReason())
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
