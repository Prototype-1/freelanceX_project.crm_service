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
	StartDate   time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate     time.Time `gorm:"type:date" json:"end_date"`
	Status      string    `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	Client Client `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
