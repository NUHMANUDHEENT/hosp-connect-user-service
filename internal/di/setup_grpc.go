package di

import (
    "log"
    "net"

    adminpb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/admin"
    docpb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/doctor"
    patientpb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/patient"
    database "github.com/nuhmanudheent/hosp-connect-user-service/internal/config"
    "github.com/nuhmanudheent/hosp-connect-user-service/internal/handler"
    "github.com/nuhmanudheent/hosp-connect-user-service/internal/repository"
    "github.com/nuhmanudheent/hosp-connect-user-service/internal/service"
    "github.com/nuhmanudheent/hosp-connect-user-service/logs"
    "google.golang.org/grpc"
    "github.com/grpc-ecosystem/go-grpc-prometheus" // Import grpc-prometheus
)

func GRPCSetup(port string) (net.Listener, *grpc.Server) {
    listener, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Failed to listen on port %s: %v", port, err)
    }
    grpc_prometheus.EnableHandlingTimeHistogram() 

    // Initialize database and repositories
    db := database.InitDatabase()
    logger := logs.NewLogger()
    adminRepo := repository.NewAdminRepository(db)
    doctorRepo := repository.NewDoctorRepository(db)
    patientRepo := repository.NewPatientRepository(db)

    // Initialize services and handlers
    adminService := service.NewAdminService(adminRepo, doctorRepo, patientRepo, logger)
    doctorService := service.NewDoctorService(doctorRepo, logger)
    patientService := service.NewPatientService(patientRepo, logger)
    adminHandler := handler.NewAdminHandler(adminService)
    doctorHandler := handler.NewDoctorHandler(doctorService)
    patientHandler := handler.NewPatientHandler(patientService)

    // Create gRPC server with Prometheus interceptors
    server := grpc.NewServer(
        grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
        grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
    )
    // Register Prometheus metrics for gRPC
    grpc_prometheus.Register(server)

    // Register services with the gRPC server
    adminpb.RegisterAdminServiceServer(server, adminHandler)
    docpb.RegisterDoctorServiceServer(server, doctorHandler)
    patientpb.RegisterPatientServiceServer(server, patientHandler)

    log.Printf("User gRPC service is running on port %s", port)
    return listener, server
}
