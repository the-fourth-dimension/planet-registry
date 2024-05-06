package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type Database struct {
	DB *gorm.DB
}

func ConnectDatabase() *Database {
	db_driver := env.GetEnv(env.DB_DRIVER)
	db_host := env.GetEnv(env.DB_HOST)
	db_user := env.GetEnv(env.DB_USER)
	db_password := env.GetEnv(env.DB_PASSWORD)
	db_name := env.GetEnv(env.DB_NAME)
	db_port := env.GetEnv(env.DB_PORT)

	db_url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", db_host, db_user, db_password, db_name, db_port)

	DB, err := gorm.Open(db_driver, db_url)

	if err != nil {
		log.Fatal("database connection error:", err)
	} else {
		fmt.Println("database connection successfull:", db_driver)
	}

	if env.GetEnv(env.APP_ENV) != "PRODUCTION" {
		DB.LogMode(true)
	}
	return &Database{DB}
}

func (d *Database) MigrateModels() {
	d.DB.AutoMigrate(&models.Admin{}, &models.Config{}, &models.Invite{}, &models.Planet{})
}

func (d *Database) PopulateConfig() {
	configRepo := repositories.NewConfigRepository(d.DB)
	var config repositories.RepositoryResult[models.Config]
	config = configRepo.FindFirst(&models.Config{})

	if config.Error != nil {
		if errors.Is(config.Error, gorm.ErrRecordNotFound) {
			log.Println("Config not found, creating one")
			config.Result = models.Config{InviteOnly: true}
			if err := configRepo.Save(&config.Result).Error; err != nil {
				log.Fatalf("Error creating config: %v", err)
			}
		} else {
			log.Fatalf("Error querying config: %v", config.Error)
		}
	}
	log.Printf("Server started with config:\n")
	spew.Dump(config.Result)
}
