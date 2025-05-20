package service

import (
	"context"
	"errors"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/repository"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/model"
	clientpb "github.com/Prototype-1/freelanceX_project.crm_service/proto/client"
"github.com/google/uuid"
"google.golang.org/grpc/metadata"
"google.golang.org/protobuf/types/known/timestamppb"
)

func parseUUID(id string) uuid.UUID {
	u, _ := uuid.Parse(id)
	return u
}

type ClientService struct {
	clientpb.UnimplementedClientServiceServer
	repo repository.ClientRepository
}

func NewClientService(repo repository.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) CreateClient(ctx context.Context, req *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	roles := md.Get("role")
	if len(roles) == 0 || roles[0] != "client" {
		return nil, errors.New("unauthorized: only clients can create client details")
	}

	    userIDs := md.Get("user_id")
    if len(userIDs) == 0 {
        return nil, errors.New("missing user ID in metadata")
    }
    
    userID, err := uuid.Parse(userIDs[0])
    if err != nil {
        return nil, errors.New("invalid user ID format")
    }

	client := &models.Client{
		ID:          userID,
		CompanyName: req.CompanyName,
		ContactName: req.ContactName,
		Email:       req.Email,
	}

	if err := s.repo.Create(ctx, client); err != nil {
		return nil, err
	}

	return &clientpb.CreateClientResponse{
		Client: convertToProto(client),
	}, nil
}

func (s *ClientService) GetClient(ctx context.Context, req *clientpb.GetClientRequest) (*clientpb.GetClientResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	roles := md.Get("role")
if len(roles) == 0 || (roles[0] != "client" && roles[0] != "admin") {
    return nil, errors.New("unauthorized: only clients or admins are allowed")
}
	client, err := s.repo.GetByID(ctx, req.ClientId)
	if err != nil {
		return nil, err
	}

	return &clientpb.GetClientResponse{
		Client: convertToProto(client),
	}, nil
}

func (s *ClientService) UpdateClient(ctx context.Context, req *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	roles := md.Get("role")
	if len(roles) == 0 || roles[0] != "client" {
		return nil, errors.New("unauthorized: only clients can update client details")
}

existing, err := s.repo.GetByID(ctx, req.ClientId)
if err != nil {
    return nil, err
}

	client := &models.Client{
		ID:          parseUUID(req.ClientId),
		CompanyName: req.CompanyName,
		ContactName: req.ContactName,
		Email:       req.Email,
		CreatedAt:   existing.CreatedAt, 
	}

	if err := s.repo.Update(ctx, client); err != nil {
		return nil, err
	}

	return &clientpb.UpdateClientResponse{
		Client: convertToProto(client),
	}, nil
}

func (s *ClientService) DeleteClient(ctx context.Context, req *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	roles := md.Get("role")
	if len(roles) == 0 || (roles[0] != "client" && roles[0] != "admin") {
		return nil, errors.New("unauthorized: only clients or admins are allowed")
}

	err := s.repo.Delete(ctx, req.ClientId)
	if err != nil {
		return nil, err
	}

	return &clientpb.DeleteClientResponse{
		Status: "deleted",
	}, nil
}

// helper to convert DB model to proto
func convertToProto(c *models.Client) *clientpb.Client {
	return &clientpb.Client{
		Id:          c.ID.String(),
		CompanyName: c.CompanyName,
		ContactName: c.ContactName,
		Email:       c.Email,
		CreatedAt:   timestamppb.New(c.CreatedAt),
	}
}
