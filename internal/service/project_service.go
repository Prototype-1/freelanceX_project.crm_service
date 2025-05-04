package service

import (
	"context"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/model"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/repository"
	projectPb "github.com/Prototype-1/freelanceX_project.crm_service/proto/project"
	"google.golang.org/protobuf/types/known/timestamppb"
	profilePb "github.com/Prototype-1/freelanceX_project.crm_service/proto/user_profile"
	"github.com/google/uuid"
	"fmt"
	"time"
	"encoding/json"
	"github.com/Prototype-1/freelanceX_project.crm_service/pkg"
)

type ProjectService struct {
	repo repository.ProjectRepository
	profileClient profilePb.ProfileServiceClient
	projectPb.UnimplementedProjectServiceServer
}

func NewProjectService(repo repository.ProjectRepository, profileClient profilePb.ProfileServiceClient) *ProjectService {
	return &ProjectService{repo: repo, profileClient: profileClient}
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

func (s *ProjectService) DiscoverProjects(ctx context.Context, req *projectPb.DiscoverProjectsRequest) (*projectPb.DiscoverProjectsResponse, error) {
	cacheKey := "discover_projects:" + req.UserId
	cached, err := pkg.Rdb.Get(pkg.Ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var cachedRes projectPb.DiscoverProjectsResponse
		err = json.Unmarshal([]byte(cached), &cachedRes)
		if err == nil {
			return &cachedRes, nil
		}
	}

	profileResp, err := s.profileClient.GetProfile(ctx, &profilePb.GetProfileRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	projects, err := s.repo.DiscoverProjects(ctx, profileResp.Skills, profileResp.Languages, profileResp.YearsOfExperience)
	if err != nil {
		return nil, err
	}

	res := &projectPb.DiscoverProjectsResponse{}
	for _, proj := range projects {
		res.Projects = append(res.Projects, &projectPb.DiscoverProject{
			ProjectId:   proj.ID.String(),
			ProjectName: proj.Title,
			Description: proj.Description,
			ClientId:    proj.ClientID.String(),
			StartDate:   timestamppb.New(proj.StartDate),
			EndDate:     timestamppb.New(proj.EndDate),
		})
	}
	jsonData, _ := json.Marshal(res)
	_ = pkg.Rdb.Set(pkg.Ctx, cacheKey, jsonData, 10*time.Minute).Err()

	return res, nil
}

func (s *ProjectService) AssignFreelancer(ctx context.Context, req *projectPb.AssignFreelancerRequest) (*projectPb.AssignFreelancerResponse, error) {
	profileResp, err := s.profileClient.GetProfile(ctx, &profilePb.GetProfileRequest{
		UserId: req.FreelancerId,
	})
	if err != nil {
		return nil, err
	}

	if len(profileResp.Skills) == 0 {
		return nil, fmt.Errorf("user is not a freelancer")
	}

	err = s.repo.AssignFreelancer(ctx, req.ProjectId, req.FreelancerId)
	if err != nil {
		return nil, err
	}

	pkg.Rdb.Del(pkg.Ctx, "discover_projects:"+req.FreelancerId)

	return &projectPb.AssignFreelancerResponse{
		ProjectId:    req.ProjectId,
		FreelancerId: req.FreelancerId,
		Status:       "assigned",
	}, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, req *projectPb.UpdateProjectRequest) (*projectPb.UpdateProjectResponse, error) {
	updateMap := map[string]interface{}{
		"title":       req.ProjectName,
		"description": req.Description,
		"end_date":    req.EndDate.AsTime(),
	}
	err := s.repo.UpdateProject(ctx, req.ProjectId, updateMap)
	if err != nil {
		return nil, err
	}

	return &projectPb.UpdateProjectResponse{
		ProjectId: req.ProjectId,
		Status:    "updated",
	}, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, req *projectPb.DeleteProjectRequest) (*projectPb.DeleteProjectResponse, error) {
	err := s.repo.DeleteProject(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}

	return &projectPb.DeleteProjectResponse{
		ProjectId: req.ProjectId,
		Status:    "deleted",
	}, nil
}

