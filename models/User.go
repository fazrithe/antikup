package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/mmattklaus/go-jwt-demo/config"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	gorm.Model
	Email    string `validate:"required,email" gorm:"default:test@mail.com"`
	Username string `validate:"required"`
	Password string `validate:"required"`
}

var validate *validator.Validate

func (u *User) FindAll(db *gorm.DB) []User {
	var users []User
	db.Find(&users)
	return users
}

func (u *User) Find(db *gorm.DB, id int) User {
	var user User
	db.First(&user, id)
	return user
}

func (u *User) Login(db *gorm.DB, username string, password string) (User, string, error) {
	validate = validator.New()
	err := validate.Var(username, "required,gt=1")
	var user User
	if err != nil {
		return user, "", errors.New("username is required")
	}
	err = validate.Var(password, "required")
	if err != nil {
		return user, "", errors.New("password not provided")
	}
	db = db.Where("username = ?", username).Find(&user)
	// .Where("password", password)
	if db.RowsAffected < 1 {
		return user, "", errors.New("invalid username")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return user, "", errors.New("wrong password provided")
	}

	token, err := (&user).GenerateToken()
	if err != nil {
		return user, "", err
	}

	return user, token, nil
}

func (u *User) Save(db *gorm.DB) (bool, error) {
	validate = validator.New()
	err := validate.Struct(u)
	if err != nil {
		return false, err
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.Password = string(hash)

	db = db.Create(&u)
	return db.RowsAffected == 1, nil
}

func (u *User) Delete(db *gorm.DB, id int) bool {
	db = db.Where("id = ?", id).Delete(User{})
	if db.RowsAffected < 1 {
		return false
	}
	return true
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	return nil
}

func (u *User) GenerateToken() (string, error) {
	var conf config.Config
	conf.Read()

	var mySigningKey = []byte(conf.AppKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = u.Username
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Printf("Error generating key: %v", err.Error())
		return "", err
	}
	return tokenString, nil
}
