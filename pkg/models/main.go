package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
)

func ConnectDatabase() *gorm.DB {

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
		fmt.Println("database connection successfull", db_driver)
	}

	DB.AutoMigrate(&Admin{}, &Config{}, &InviteCode{}, &Planet{})

	return DB
}
