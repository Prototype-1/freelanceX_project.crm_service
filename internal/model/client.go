package models

import (
	"time"
	"github.com/google/uuid"
)

type Client struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CompanyName  string    `gorm:"type:varchar(255);not null" json:"company_name"`
	ContactName  string    `gorm:"type:varchar(255)" json:"contact_name"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
