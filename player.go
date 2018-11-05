package main

import (
	"encoding/json"
	"net/http"

	"github.com/garyburd/redigo/redis"
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

	if vars["name"] == "" {
		sendError(w, 400, "Missing parameter: name")
		return
	}

	conn := kt.redis.Get()
	defer conn.Close()

	result, err := redis.String(conn.Do("GET", objectPrefix+vars["roomId"]))

	if err != nil {
		sendError(w, 500, err.Error())
		return
	}

	var room Room
	err = json.Unmarshal([]byte(result), &room)

	if err != nil {
		sendError(w, 500, err.Error())
		return
	}

	player := Player{vars["name"], 0}
	room.Players = append(room.Players, player)

	updated, err := updateRoom(room)
	sendResp(w, updated)
}
