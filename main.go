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
	kt     *Kaladont
	router *mux.Router
)

// our main function
func main() {
	fmt.Println("Starting program")
	kt = initKaladont()
	router = mux.NewRouter()
	AuthRouter()

	router.HandleFunc("/game/{roomID}", GetGame).Methods("GET")
	router.HandleFunc("/create", CreateGame).Queries("name", "{playerName}").Methods("POST")
	router.HandleFunc("/join/{roomId}", isAuthenticated(JoinRoom)).Queries("name", "{name}").Methods("POST")
	router.HandleFunc("/{roomId}", RemoveGame).Methods("DELETE")
	router.HandleFunc("/{roomId}/leave", isAuthenticated(LeaveRoom)).Methods("DELETE")

	// router.
	log.Fatal(http.ListenAndServe(":6969", handlers.LoggingHandler(os.Stdout, router)))
}
