package handler

import (
	"context"
	"encoding/json"
	"fmt"
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
	fmt.Println("pat", req.PatientId)
	patient, err := p.service.GetProfile(req.PatientId)
	if err != nil {
		return &pb.GetProfileResponse{
			Status:     "fail",
			StatusCode: 404,
			Error:      err.Error(),
		}, nil
	}

	fmt.Println("e", patient.Email)
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
func (p *PatientServiceClient) AddPrescription(ctx context.Context, req *pb.AddPrescriptionRequest) (*pb.StandardResponse, error) {
	log.Printf("request for storing prescription from doctorid : %s to patient id: %s", req.DoctorId, req.PatientId)
	var patientPrescription domain.PatientPrescription

	patientPrescription.PatientId = req.PatientId
	patientPrescription.DoctorId = req.DoctorId
	var presc []domain.Prescription
	for _, v := range req.Prescription {
		prescription := domain.Prescription{
			Medication: v.Medication,
			Dosage:     v.Dosage,
			Frequency:  v.Frequency,
		}
		presc = append(presc, prescription)
	}
	jsonPrescriptions, err := json.Marshal(presc)
	if err != nil {
		return &pb.StandardResponse{
			Status:     "fail",
			Message:    err.Error(),
			StatusCode: 400,
		}, nil
	}
	patientPrescription.Prescription = string(jsonPrescriptions)
	// Store the patientPrescription object into the database
	if err := p.service.AddPrescription(patientPrescription); err != nil {
		return &pb.StandardResponse{
			Status:     "fail",
			Message:    err.Error(),
			StatusCode: 400,
		}, nil
	}
	// Return a success response
	return &pb.StandardResponse{
		Status:  "success",
		Message: "Prescription added successfully",
	}, nil
}
func (p PatientServiceClient) GetPrescription(ctx context.Context, req *pb.GetPrescriptionRequest) (*pb.GetPrescriptionResponse, error) {
	log.Printf("request for getting prescription from patient id: %s  , %s", req.PatientId, req.Query)
	resp, err := p.service.GetPrescription(req.PatientId, req.Query)
	if err != nil {
		return &pb.GetPrescriptionResponse{
			Prescriptions: nil,
			Status:        "fail",
			StatusCode:    "400",
			Message:       "No prescription found",
		}, nil
	}
	var prescriptions []*pb.GetPresc
	for _, v := range resp {
		var presc []*pb.Prescription
		if err := json.Unmarshal([]byte(v.Prescription), &presc); err != nil {
			log.Printf("error unmarshalling prescription for patient %s: %v", req.PatientId, err)
			return &pb.GetPrescriptionResponse{
				Prescriptions: nil,
				Status:        "fail",
				StatusCode:    "500",
				Message:       "Failed to process prescriptions",
			}, nil
		}
		prescriptions = append(prescriptions, &pb.GetPresc{
			DoctorId:         v.DoctorId,
			PrescriptionTime: v.CreatedAt.Format("2006-01-02"),
			Treatment:        presc,
		})

	}
	fmt.Println("pres", prescriptions)

	return &pb.GetPrescriptionResponse{
		Prescriptions: prescriptions,
		Status:        "success",
		StatusCode:    "200",
	}, nil
}
