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
	UpdateProject(ctx context.Context, projectID string, updated map[string]interface{}) error
DeleteProject(ctx context.Context, projectID string) error
AssignFreelancer(ctx context.Context, projectID, freelancerID string) error
DiscoverProjects(ctx context.Context, skills []string, languages []string, experienceMin int32) ([]models.Project, error)

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

func (r *projectRepository) UpdateProject(ctx context.Context, projectID string, updated map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&models.Project{}).Where("id = ?", projectID).Updates(updated).Error
}

func (r *projectRepository) DeleteProject(ctx context.Context, projectID string) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, "id = ?", projectID).Error
}

func (r *projectRepository) AssignFreelancer(ctx context.Context, projectID, freelancerID string) error {
	return r.db.WithContext(ctx).Exec("INSERT INTO project_freelancers (project_id, freelancer_id) VALUES (?, ?)", projectID, freelancerID).Error
}

func (r *projectRepository) DiscoverProjects(ctx context.Context, skills []string, languages []string, experienceMin int32) ([]models.Project, error) {
	var projects []models.Project

	err := r.db.WithContext(ctx).
		Where("status = ?", "ongoing").
		Where("required_skills && ?", skills).
		Where("required_languages && ?", languages).
		Where("min_experience <= ?", experienceMin). 
		Find(&projects).Error

	return projects, err
}
