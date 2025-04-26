package repository

import (
	"github.com/google/uuid"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/model"
	"gorm.io/gorm"
	"errors"
)

type Filters struct {
	Skills           []string
	ExperienceMin    int
	Languages        []string
	ProjectStatus    string
	AssignedFreelancer uuid.UUID
}

type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetProjectsByUserID(userID uuid.UUID) ([]models.Project, error)
	GetProjectByID(projectID uuid.UUID) (*models.Project, error)
	DiscoverProjects(filters *Filters) ([]models.Project, error)
	AssignFreelancer(projectID, freelancerID uuid.UUID) error
	UpdateProject(project *models.Project) error
	DeleteProject(projectID uuid.UUID) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) CreateProject(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) GetProjectsByUserID(userID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Where("client_id = ? OR id IN (?)", userID, r.getAssignedProjects(userID)).Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) GetProjectByID(projectID uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Client").First(&project, "id = ?", projectID).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) DiscoverProjects(filters *Filters) ([]models.Project, error) {
	var projects []models.Project
	query := r.db.Model(&models.Project{}).Preload("Client")

	if filters.Skills != nil {
		query = query.Where("skills @> ?", filters.Skills)
	}
	if filters.ExperienceMin > 0 {
		query = query.Where("experience >= ?", filters.ExperienceMin)
	}
	if filters.Languages != nil {
		query = query.Where("languages @> ?", filters.Languages) 
	}
	if filters.ProjectStatus != "" {
		query = query.Where("status = ?", filters.ProjectStatus)
	}
	if filters.AssignedFreelancer != uuid.Nil {
		query = query.Where("id NOT IN (?)", r.getAssignedProjects(filters.AssignedFreelancer))
	}

	err := query.Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) getAssignedProjects(freelancerID uuid.UUID) []uuid.UUID {
	var projectIDs []uuid.UUID
	r.db.Table("projects").
		Joins("INNER JOIN freelancer_assignments ON freelancer_assignments.project_id = projects.id").
		Where("freelancer_assignments.freelancer_id = ?", freelancerID).
		Pluck("projects.id", &projectIDs)
	return projectIDs
}

// AssignFreelancer assigns a freelancer to a project
func (r *projectRepository) AssignFreelancer(projectID, freelancerID uuid.UUID) error {
	// Assuming we have a table `freelancer_assignments` to track freelancer assignments
	// Add check for existing assignment before inserting
	exists := r.db.Table("freelancer_assignments").Where("project_id = ? AND freelancer_id = ?", projectID, freelancerID).Exists()
	if exists {
		return errors.New("freelancer already assigned")
	}

	// Assign freelancer to project
	return r.db.Table("freelancer_assignments").Create(&FreelancerAssignment{
		ProjectID:  projectID,
		FreelancerID: freelancerID,
	}).Error
}

func (r *projectRepository) UpdateProject(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) DeleteProject(projectID uuid.UUID) error {
	return r.db.Delete(&models.Project{}, "id = ?", projectID).Error
}