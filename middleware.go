package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func parseRoomFromParams(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var roomID = mux.Vars(r)["roomId"]
		room, err := getRoom(roomID)

		if err != nil {
			sendError(w, http.StatusNotFound, "Room not found")
			return
		}

		var _room Room
		json.Unmarshal([]byte(room), &_room)

		context.Set(r, "room", _room)
		next.ServeHTTP(w, r)
	}
}
