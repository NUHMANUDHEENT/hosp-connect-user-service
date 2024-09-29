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

func (a *AdminServiceClient) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	admindetails := domain.Admin{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	log.Println("signin request", admindetails)
	resp, err := a.service.SignIn(admindetails)
	if err != nil {
		return &pb.SignInResponse{
			Status:     "fail",
			Message:    "Invalid credentials, please try again.",
			StatusCode: 401,
		}, nil

	}
	return &pb.SignInResponse{
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
func (a *AdminServiceClient) AddSpecialization(ctx context.Context, req *pb.AddSpecializationRequest) (*pb.StandardResponse, error) {
	log.Println("Adding specialization with name: ", req.Name)
	resp, err := a.service.AddSpecialization(req.Name, req.Description)
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
func (a *AdminServiceClient) ListDoctors(ctx context.Context, req *pb.Empty) (*pb.ListDoctorsResponse, error) {
	doctors, err := a.service.ListDoctors()
	if err != nil {
		return &pb.ListDoctorsResponse{
			Status:     "fail",
			StatusCode: 500,
			Error:      err.Error(),
		}, nil
	}

	var pbDoctors []*pb.Doctor
	for _, doctor := range doctors {
		pbDoctors = append(pbDoctors, &pb.Doctor{
			DoctorId:         doctor.DoctorId,
			Name:             doctor.Name,
			Email:            doctor.Email,
			Phone:            int32(doctor.Phone),
			SpecializationId: int32(doctor.SpecilazationId),
		})
	}

	return &pb.ListDoctorsResponse{
		Doctors:    pbDoctors,
		Status:     "success",
		StatusCode: 200,
	}, nil
}

// ListPatients returns a list of all patients
func (a *AdminServiceClient) ListPatients(ctx context.Context, req *pb.Empty) (*pb.ListPatientsResponse, error) {
	patients, err := a.service.ListPatients()
	if err != nil {
		return &pb.ListPatientsResponse{
			Status:     "fail",
			StatusCode: 500,
			Error:      err.Error(),
		}, nil
	}

	var pbPatients []*pb.Patient
	for _, patient := range patients {
		pbPatients = append(pbPatients, &pb.Patient{
			PatientId: patient.PatientID,
			Name:      patient.Name,
			Email:     patient.Email,
			Phone:     int32(patient.Phone),
			Age:       int32(patient.Age),
			Gender:    patient.Gender,
		})
	}

	return &pb.ListPatientsResponse{
		Patients:   pbPatients,
		Status:     "success",
		StatusCode: 200,
	}, nil
}
