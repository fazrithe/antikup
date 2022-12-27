package models

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fazrithe/antikup/config"
	"github.com/jinzhu/gorm"
	"github.com/mmattklaus/go-jwt-demo/helpers"
)

type Config2 struct {
	Db  *gorm.DB
	Env *config.Config
	Log *log.Logger
}

var (
	user *User
)

func (c *Config2) Index(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithJson(w, 200, user.FindAll(c.Db))
}

func (c *Config2) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error reading request data: %v", err.Error()))
		return
	}

	u, token, err := user.Login(c.Db, r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.RespondWithJson(w, 200, map[string]interface{}{"user": u, "token": token})
}

func (c *Config2) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error reading request data: %v", err.Error()))
		return
	}
	newUser := User{}
	newUser.Email = r.FormValue("email")
	newUser.Password = r.FormValue("password")
	newUser.Username = r.FormValue("username")

	_, err := newUser.Save(c.Db)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Error while creating User. %v", err.Error()))
		return
	}
	helpers.RespondWithJson(w, 200, map[string]string{"status": "success", "message": "User created."})
}

func (c *Config2) Delete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error reading request data: %v", err.Error()))
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Bad data provided: %v", err.Error()))
		return
	}
	ok := user.Delete(c.Db, id)
	if !ok {
		helpers.RespondWithError(w, http.StatusNotFound, "not found")
		return
	}
	helpers.RespondWithJson(w, http.StatusOK, map[string]string{"status": "ok", "message": "user deleted"})
}

func (c *Config2) Find(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error reading request data: %v", err.Error()))
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Bad data provided: %v", err.Error()))
		return
	}
	u := user.Find(c.Db, id)
	if u.ID > 0 {
		helpers.RespondWithJson(w, http.StatusOK, u)
		return
	}
	helpers.RespondWithError(w, http.StatusNotFound, "not found")
}

func UserRoutes(db *gorm.DB, conf *config.Config, logger *log.Logger) *Config2 {
	return &Config2{
		Db:  db,
		Env: conf,
		Log: logger,
	}
}
