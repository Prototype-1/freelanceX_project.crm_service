package repository

import (
	"context"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/model"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, project *models.Project) error
	GetProjectsByClientID(ctx context.Context, clientID string) ([]models.Project, error)
	GetProjectByID(ctx context.Context, projectID string) (*models.Project, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) CreateProject(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectRepository) GetProjectsByClientID(ctx context.Context, clientID string) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.WithContext(ctx).Where("client_id = ?", clientID).Find(&projects).Error
	return projects, err
}

func (r *projectRepository) GetProjectByID(ctx context.Context, projectID string) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).First(&project, "id = ?", projectID).Error
	return &project, err
}