package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mmattklaus/go-jwt-demo/config"
)

type Database struct {
	DB *gorm.DB
}

func (c *Database) Connect(config *config.Config, logger *log.Logger) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", config.DbUsername, config.DbPassword, config.DbDatabase))
	if err != nil {
		logger.Panicf("Database connection error: %v", err.Error())
	}
	db.SetLogger(logger)
	c.DB = db
	logger.Println("Database connection established.")
	// defer db.Close()
}
