package service

import (
	"context"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/model"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/repository"
	projectPb "github.com/Prototype-1/freelanceX_project.crm_service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/google/uuid"
)

type ProjectService struct {
	repo repository.ProjectRepository
	projectPb.UnimplementedProjectServiceServer
}

func NewProjectService(repo repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(ctx context.Context, req *projectPb.CreateProjectRequest) (*projectPb.CreateProjectResponse, error) {
	id := uuid.New()
	project := &models.Project{
		ID:          id,
		ClientID:    uuid.MustParse(req.ClientId),
		Title:       req.ProjectName,
		Description: req.Description,
		StartDate:   req.StartDate.AsTime(),
		EndDate:     req.EndDate.AsTime(),
		Status:      "ongoing",
	}
	err := s.repo.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return &projectPb.CreateProjectResponse{
		ProjectId: id.String(),
		Status:    "created",
	}, nil
}

func (s *ProjectService) GetProjectsByUser(ctx context.Context, req *projectPb.GetProjectsByUserRequest) (*projectPb.GetProjectsByUserResponse, error) {
	projects, err := s.repo.GetProjectsByClientID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	res := &projectPb.GetProjectsByUserResponse{}
	for _, proj := range projects {
		res.Projects = append(res.Projects, &projectPb.ProjectSummary{
			ProjectId:  proj.ID.String(),
			ProjectName: proj.Title,
			Role:       "client",
			StartDate:  timestamppb.New(proj.StartDate),
			EndDate:    timestamppb.New(proj.EndDate),
		})
	}

	return res, nil
}

func (s *ProjectService) GetProjectById(ctx context.Context, req *projectPb.GetProjectByIdRequest) (*projectPb.GetProjectByIdResponse, error) {
	project, err := s.repo.GetProjectByID(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}

	res := &projectPb.GetProjectByIdResponse{
		ProjectId:   project.ID.String(),
		ProjectName: project.Title,
		Description: project.Description,
		ClientId:    project.ClientID.String(),
		StartDate:   timestamppb.New(project.StartDate),
		EndDate:     timestamppb.New(project.EndDate),
		AssignedFreelancers: []*projectPb.FreelancerInfo{}, 
	}

	return res, nil
}