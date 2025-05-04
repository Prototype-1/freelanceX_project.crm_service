package main

import (
	"log"
	"net"
	"os"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	projectPb "github.com/Prototype-1/freelanceX_project.crm_service/proto/project" 
	clientpb "github.com/Prototype-1/freelanceX_project.crm_service/proto/client" 
	"github.com/Prototype-1/freelanceX_project.crm_service/migrations"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/repository"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/service"
	"github.com/Prototype-1/freelanceX_project.crm_service/client"
	redisDb "github.com/Prototype-1/freelanceX_project.crm_service/pkg"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on system environment")
	}
}

func main() {
	db := migrations.ConnectDatabase()
	migrations.RunMigrations(db)

	redisDb.InitRedis()

	clientRepo := repository.NewClientRepository(db)
clientSvc := service.NewClientService(clientRepo)

	projectRepo := repository.NewProjectRepository(db)
	profileClient := client.NewProfileServiceClient() 
	projectService := service.NewProjectService(projectRepo, profileClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "50053"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	projectPb.RegisterProjectServiceServer(grpcServer, projectService)
	clientpb.RegisterClientServiceServer(grpcServer, clientSvc)

	log.Printf("gRPC server is running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
