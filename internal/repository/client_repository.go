package repository

import (
	"context"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/model"
	"gorm.io/gorm"
)

type ClientRepository interface {
	Create(ctx context.Context, client *models.Client) error
	GetByID(ctx context.Context, id string) (*models.Client, error)
	Update(ctx context.Context, client *models.Client) error
	Delete(ctx context.Context, id string) error
}

type clientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(ctx context.Context, client *models.Client) error {
	return r.db.WithContext(ctx).Create(client).Error
}

func (r *clientRepository) GetByID(ctx context.Context, id string) (*models.Client, error) {
	var client models.Client
	err := r.db.WithContext(ctx).First(&client, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) Update(ctx context.Context, client *models.Client) error {
	return r.db.WithContext(ctx).Save(client).Error
}

func (r *clientRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Client{}, "id = ?", id).Error
}
