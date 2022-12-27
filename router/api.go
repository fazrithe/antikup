package router

import (
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mmattklaus/go-jwt-demo/config"
	mw "github.com/mmattklaus/go-jwt-demo/middleware"
	"github.com/mmattklaus/go-jwt-demo/models"
)

type Config2 struct {
	Db  *gorm.DB
	Env *config.Config
	Log *log.Logger
}

func (c *Config2) InitializeRoutes() {
	User := models.UserRoutes(c.Db, c.Env, c.Log)
	http.HandleFunc("/register", c.Logger(User.Register))
	http.HandleFunc("/login", c.Logger(User.Login))
	http.HandleFunc("/users", mw.IsAuthorized(c.Logger(User.Index)))
	http.HandleFunc("/delete", mw.IsAuthorized(c.Logger(User.Delete)))
	http.HandleFunc("/find", mw.IsAuthorized(c.Logger(User.Find)))
}

func NewAPI(db *gorm.DB, conf *config.Config, logger *log.Logger) *Config2 {
	return &Config2{
		Db:  db,
		Env: conf,
		Log: logger,
	}
}
func (c *Config2) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer c.Log.Printf("%v %v %v %v in %v", r.Proto, r.Method, r.Host, r.RequestURI, time.Since(start))
		next(w, r)
	}
}
