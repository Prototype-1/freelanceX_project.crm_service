package main

import (
	"log"
	"github.com/joho/godotenv"
	"github.com/Prototype-1/freelanceX_project.crm_service/migrations" 
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on system environment")
	}
}

func main() {
	db := migrations.ConnectDatabase()

	migrations.RunMigrations(db)

}
