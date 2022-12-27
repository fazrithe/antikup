package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mmattklaus/go-jwt-demo/config"
	"github.com/mmattklaus/go-jwt-demo/database"
	"github.com/mmattklaus/go-jwt-demo/models"
	"github.com/mmattklaus/go-jwt-demo/router"
)

var dh database.Database
var logger *log.Logger
var conf config.Config

var (
	User models.User
)

func main() {

	logger = log.New(os.Stdout, "Micros: ", log.LstdFlags|log.Lshortfile)
	conf.Read()
	dh.Connect(&conf, logger)
	defer dh.DB.Close()

	database.InitMigrations(dh.DB, &conf)

	api := router.NewAPI(dh.DB, &conf, logger)
	api.InitializeRoutes() // Run migrations. This will only create tables & fields which don't exist.

	logger.Println("server started on PORT: " + conf.ServerAddr)

	logger.Fatalln(http.ListenAndServe(conf.ServerAddr, nil))

}
