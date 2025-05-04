package service

import (
	"context"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/repository"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/model"
	clientpb "github.com/Prototype-1/freelanceX_project.crm_service/proto/client"
"github.com/google/uuid"
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
	client := &models.Client{
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
	client, err := s.repo.GetByID(ctx, req.ClientId)
	if err != nil {
		return nil, err
	}

	return &clientpb.GetClientResponse{
		Client: convertToProto(client),
	}, nil
}

func (s *ClientService) UpdateClient(ctx context.Context, req *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error) {
	client := &models.Client{
		ID:          parseUUID(req.ClientId),
		CompanyName: req.CompanyName,
		ContactName: req.ContactName,
		Email:       req.Email,
	}

	if err := s.repo.Update(ctx, client); err != nil {
		return nil, err
	}

	return &clientpb.UpdateClientResponse{
		Client: convertToProto(client),
	}, nil
}

func (s *ClientService) DeleteClient(ctx context.Context, req *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error) {
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
