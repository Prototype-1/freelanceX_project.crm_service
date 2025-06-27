package models

import (
	"time"
	"github.com/google/uuid"
)

type ProjectFreelancer struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProjectID    uuid.UUID `gorm:"type:uuid;not null" json:"project_id"`
	FreelancerID uuid.UUID `gorm:"type:uuid;not null" json:"freelancer_id"`
	AssignedAt   time.Time `gorm:"autoCreateTime" json:"assigned_at"`
	ProjectIDRef Project   `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

