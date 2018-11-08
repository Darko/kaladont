// NaM
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

var (
	router *mux.Router
	db     Db
)

// our main function
func main() {
	fmt.Println("Starting program")
	db = Db{Conn: initRedis().Get()}
	router = mux.NewRouter()

	AuthRouter()

	router.HandleFunc("/v1/games", CreateGame).Methods("POST")
	router.HandleFunc("/v1/games/{roomId}", GetGame).Methods("GET")
	router.HandleFunc("/v1/games/{roomId}", RemoveGame).Methods("DELETE")
	router.HandleFunc("/v1/games/{roomId}/join", isAuthenticated(JoinRoom)).Methods("POST")
	router.HandleFunc("/v1/games/{roomId}/leave", isAuthenticated(LeaveRoom)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":6969", handlers.LoggingHandler(os.Stdout, router)))
}
