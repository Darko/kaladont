package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

var jwtSecret = []byte("topsecret")

// AuthRouter creates a subrouter of the main router to handle auth requests
func AuthRouter() *mux.Router {
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", Authenticate).Methods("POST")
	return authRouter
}

// Authenticate authenticates the user upon an http request
// given theres a name provided, if there's not, an error is returned
func Authenticate(w http.ResponseWriter, r *http.Request) {
	var player Player
	err := parseBody(r, &player)

	if err != nil {
		sendError(w, 500, err.Error())
		return
	}

	if player.Name == "" {
		sendError(w, 400, "Missing property: name")
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = player.Name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, _ := token.SignedString(jwtSecret)
	cookie := http.Cookie{
		Name:    "kaladont-token",
		Value:   tokenString,
		Expires: time.Now().AddDate(0, 0, 30),
	}

	// w.Header().Set("Authorization", "Bearer "+tokenString)
	http.SetCookie(w, &cookie)
	sendResp(w, map[string]string{
		"token": tokenString,
	})
}

func isAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authHeader = r.Header.Get("Authorization")
		var t = strings.Split(authHeader, " ")

		if len(t) < 2 {
			sendError(w, 401, "Invalid token")
			return
		}

		token, _ := jwt.Parse(t[1], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return jwtSecret, nil
		})

		if val, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			context.Set(r, "tokenData", val)
			next.ServeHTTP(w, r)
		} else {
			sendError(w, 401, "Invalid token")
		}
	}
}
