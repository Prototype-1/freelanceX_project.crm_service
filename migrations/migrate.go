package migrations

import (
	"fmt"
	"log"
	"os"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/Prototype-1/freelanceX_project.crm_service/internal/model" 
)

func ConnectDatabase() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db
}

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Client{},
		&models.Project{},
		&models.ProjectFreelancer{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migrated successfully!")
}
