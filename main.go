// NaM
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	kt *Kaladont
)

// our main function
func main() {
	fmt.Println("Starting program")
	kt = initKaladont()
	router := mux.NewRouter()
	router.HandleFunc("/game/{roomID}", GetGame).Methods("GET")
	router.HandleFunc("/create", CreateGame).Queries("name", "{playerName}").Methods("POST")
	router.HandleFunc("/join/{roomId}", JoinRoom).Queries("name", "{name}").Methods("POST")
	router.HandleFunc("/{roomId}", RemoveGame).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":6969", router))
}
