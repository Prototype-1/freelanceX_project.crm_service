package models

import (
	"time"
	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ClientID    uuid.UUID `gorm:"type:uuid;not null" json:"client_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	RequiredSkills    []string `gorm:"type:text[]"` 
	MinExperience int32 `gorm:"type:integer" json:"min_experience"` 
RequiredLanguages []string `gorm:"type:text[]"`
	StartDate   time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate     time.Time `gorm:"type:date" json:"end_date"`
	Freelancers []ProjectFreelancer `gorm:"foreignKey:ProjectID" json:"-"`
	Status      string    `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	Client Client `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
