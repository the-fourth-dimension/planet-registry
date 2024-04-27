package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
)

var DB *gorm.DB

func ConnectDatabase() {

	Dbdriver := env.GetEnv(env.DB_DRIVER)
	DbHost := env.GetEnv(env.DB_HOST)
	DbUser := env.GetEnv(env.DB_USER)
	DbPassword := env.GetEnv(env.DB_PASSWORD)
	DbName := env.GetEnv(env.DB_NAME)
	DbPort := env.GetEnv(env.DB_PORT)

	DBURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DbHost, DbUser, DbPassword, DbName, DbPort)

	DB, err := gorm.Open(Dbdriver, DBURL)

	if err != nil {
		log.Fatal("database connection error:", err)
	} else {
		fmt.Println("database connection successfull", Dbdriver)
	}

	DB.AutoMigrate()
}
