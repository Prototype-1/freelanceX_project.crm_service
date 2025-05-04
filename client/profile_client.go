package client

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	profilePb "github.com/Prototype-1/freelanceX_project.crm_service/proto/user_profile"
)

func NewProfileServiceClient() profilePb.ProfileServiceClient {
	conn, err := grpc.NewClient("user_service_host:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Could not connect to user service: %v", err)
	}
	return profilePb.NewProfileServiceClient(conn)
}
