package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/garyburd/redigo/redis"
)

const objectPrefix = "kaladont:"

// Room struct
type Room struct {
	ID      string   `json:"id"`
	Players []Player `json:"players,omitempty"`
}

// CreateGame Controller
func CreateGame(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["playerName"]
	creator := Player{name, 0}
	room := Room{randomID(5), []Player{creator}}
	conn := kt.redis.Get()
	defer conn.Close()

	b, err := json.Marshal(&room)
	if err != nil {
		sendError(w, 500, "Something got bamboozled with redis while creating a game")
		return
	}

	_, err = conn.Do("SET", "kaladont:"+room.ID, string(b))
	if err != nil {
		sendError(w, 500, "Redis command failed: "+err.Error())
		return
	}

	sendResp(w, map[string]interface{}{
		"message":    "Successfully created game",
		"statusCode": 200,
		"room":       room,
	})
}

// GetGame controller
func GetGame(w http.ResponseWriter, r *http.Request) {
	var roomID = mux.Vars(r)["roomID"]
	conn := kt.redis.Get()
	defer conn.Close()
	gameStr, err := redis.String(conn.Do("GET", objectPrefix+roomID))

	fmt.Println(gameStr)

	if err != nil {
		fmt.Println(err.Error())
		sendError(w, 500, err.Error())
		return
	}

	b := []byte(gameStr)
	var game Room
	err = json.Unmarshal(b, &game)

	if err != nil {
		sendError(w, 500, err.Error())
		return
	}

	sendResp(w, game)
}
