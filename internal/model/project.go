package models

import (
	"time"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Project struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ClientID          uuid.UUID      `gorm:"type:uuid;not null" json:"client_id"`
	Title             string         `gorm:"type:varchar(255);not null" json:"title"`
	Description       string         `gorm:"type:text" json:"description"`
	RequiredSkills    pq.StringArray `gorm:"type:text[]" json:"required_skills"`
	MinExperience     int32          `gorm:"type:integer" json:"min_experience"`
	RequiredLanguages pq.StringArray `gorm:"type:text[]" json:"required_languages"`
	StartDate         time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate           time.Time      `gorm:"type:date" json:"end_date"`
	Freelancers       []ProjectFreelancer `gorm:"foreignKey:ProjectID" json:"-"`
	Status            string         `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`

	Client Client `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}