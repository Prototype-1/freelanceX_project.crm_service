package client

import (
    "log"
    "os"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    profilePb "github.com/Prototype-1/freelanceX_project.crm_service/proto/user_profile"
)

func NewProfileServiceClient() profilePb.ProfileServiceClient {
  
    userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
    if userServiceAddr == "" {
        userServiceAddr = "localhost:50051" 
    }
    
    log.Printf("Connecting to user service at: %s", userServiceAddr)
    
    conn, err := grpc.NewClient(userServiceAddr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("Could not connect to user service at %s: %v", userServiceAddr, err)
    }
    
    log.Printf("Successfully connected to user service at: %s", userServiceAddr)
    return profilePb.NewProfileServiceClient(conn)
}