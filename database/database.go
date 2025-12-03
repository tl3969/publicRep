package database

import (
	"fmt"

	"golang_system/config"
	"golang_system/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	cfg := config.Load()
	dbConfig := cfg.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移
	err = DB.AutoMigrate(&models.User{}, &models.Article{}, &models.Comment{})
	if err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}
