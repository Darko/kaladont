package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Player struct
type Player struct {
	Name  string `json:"name"`
	Score int32  `json:"score"`
}

// JoinRoom controller
func JoinRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	// json.NewEncoder(w).Encode(room)
}
