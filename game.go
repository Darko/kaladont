package main

import (
	"encoding/json"
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

func updateRoom(room Room) (interface{}, error) {
	b, err := json.Marshal(&room)
	if err != nil {
		return nil, err
	}

	// _, err = kt.redis.Do("SET", "kaladont:"+room.ID, string(b))
	_, err = db.Set(objectPrefix+room.ID, string(b))
	return room, err
}

func removeRoom(roomID string) error {
	_, err := db.Delete(objectPrefix + roomID)
	return err
}

func getRoom(roomID string) (string, error) {
	room, err := redis.String(db.Get(objectPrefix + roomID))
	return room, err
}

// CreateGame Controller
func CreateGame(w http.ResponseWriter, r *http.Request) {
	var creator Player
	err := parseBody(r, &creator)
	room := Room{randomID(5), []Player{creator}}

	b, err := json.Marshal(&room)
	if err != nil {
		sendError(w, 500, "Something got bamboozled with redis while creating a game")
		return
	}

	_, err = db.Set(objectPrefix+room.ID, string(b))
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
	var roomID = mux.Vars(r)["roomId"]
	gameStr, err := redis.String(db.Get(objectPrefix + roomID))

	if err != nil {
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

// RemoveGame removes a game room from redis
func RemoveGame(w http.ResponseWriter, r *http.Request) {
	var roomID = mux.Vars(r)["roomId"]

	if roomID == "" {
		sendError(w, 400, "Missing parameter roomId")
		return
	}

	err := removeRoom(roomID)

	if err != nil {
		sendError(w, 500, "Server shit itself")
		return
	}

	sendResp(w, map[string]interface{}{
		"statusCode": 201,
	})
}
