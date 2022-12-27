package database

import (
	"github.com/jinzhu/gorm"
	"github.com/mmattklaus/go-jwt-demo/config"
	"github.com/mmattklaus/go-jwt-demo/models"
)

func init() {

}

func InitMigrations(db *gorm.DB, conf *config.Config) {
	db.AutoMigrate(&models.User{}) // Add any models to this list (pointer to model)
}
