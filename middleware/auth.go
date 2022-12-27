package middleware

import (
	"errors"
	"fmt"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/mmattklaus/go-jwt-demo/config"
	"github.com/mmattklaus/go-jwt-demo/helpers"
)

func IsAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (i interface{}, e error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unusual identity")
				}
				var conf config.Config
				conf.Read()
				return []byte(conf.AppKey), nil
			})
			if err != nil {
				helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				fmt.Printf(" User: %s \n ID: %v\n Exp: %v \n", claims["user"], claims["id"], claims["exp"])
				next(w, r)
				return
			}
		}
		helpers.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
}
