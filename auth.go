package main

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var jwtSecret = []byte("topsecret")

func AuthRouter() *mux.Router {
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", Authenticate).Methods("GET")
	return authRouter
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	reqBody := r.Body

	fmt.Println(reqBody)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Darko"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, _ := token.SignedString(jwtSecret)

	w.Header().Set("Authorization", "Bearer "+tokenString)
	sendResp(w, map[string]string{
		"token": tokenString,
	})
}
