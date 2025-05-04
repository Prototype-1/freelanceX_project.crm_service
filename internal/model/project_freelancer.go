package models

import (
	"time"
	"github.com/google/uuid"
)

type ProjectFreelancer struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID    uuid.UUID `gorm:"type:uuid;not null"`
	FreelancerID uuid.UUID `gorm:"type:uuid;not null"`
	AssignedAt   time.Time `gorm:"autoCreateTime"`

	ProjectIDRef Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}
