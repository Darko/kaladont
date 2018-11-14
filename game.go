package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/garyburd/redigo/redis"
)

// updateRoom JSON.Marshal's a Room type and saves it to redis
func updateRoom(room Room) (interface{}, error) {
	b, err := json.Marshal(&room)
	if err != nil {
		return nil, err
	}

	// _, err = kt.redis.Do("SET", "kaladont:"+room.ID, string(b))
	_, err = db.Set(objectPrefix+room.ID, string(b))
	return room, err
}

// removeRoom completely removes a room from redis
func removeRoom(roomID string) error {
	_, err := db.Delete(objectPrefix + roomID)
	return err
}

// getRoom returns a stringified value of Room. You must json.Unmarshal it yourself
func getRoom(roomID string) (string, error) {
	room, err := redis.String(db.Get(objectPrefix + roomID))
	return room, err
}

func serveNextPlayer(room *Room) Player {
	player, _, _ := findPlayer(room.Players, room.CurrentPlayer.Name)
	room.Players = removePlayerByName(room.Players, room.CurrentPlayer.Name)
	room.Players = append(room.Players, player)
	return room.Players[0]
}

// CreateGame Controller
func CreateGame(w http.ResponseWriter, r *http.Request) {
	var creator Player
	err := parseBody(r, &creator)
	room := Room{
		ID:            randomID(5),
		Players:       []Player{creator},
		Words:         []Word{},
		CurrentPlayer: creator,
	}

	b, err := json.Marshal(&room)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Something got bamboozled with redis while creating a game")
		return
	}

	_, err = db.Set(objectPrefix+room.ID, string(b))
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Redis command failed: "+err.Error())
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
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	b := []byte(gameStr)
	var game Room
	err = json.Unmarshal(b, &game)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResp(w, game)
}

// RemoveGame removes a game room from redis
func RemoveGame(w http.ResponseWriter, r *http.Request) {
	var roomID = mux.Vars(r)["roomId"]

	if roomID == "" {
		sendError(w, http.StatusBadRequest, "Missing parameter roomId")
		return
	}

	err := removeRoom(roomID)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Server shit itself")
		return
	}

	sendResp(w, map[string]interface{}{
		"statusCode": 201,
	})
}

// SubmitWord appends a word to the Words slice of Room
func SubmitWord(w http.ResponseWriter, r *http.Request) {
	var roomID = mux.Vars(r)["roomId"]

	if roomID == "" {
		sendError(w, http.StatusBadRequest, "Invalid room")
		return
	}

	room, err := getRoom(roomID)

	if err != nil {
		sendError(w, http.StatusNotFound, "Room not found")
		return
	}

	var _room Room
	var word Word
	json.Unmarshal([]byte(room), &_room)
	parseBody(r, &word)

	if word.Word == "" {
		sendError(w, http.StatusBadRequest, "Invalid word")
		return
	}

	if word.Player == "" {
		sendError(w, http.StatusBadRequest, "Invalid player")
		return
	}

	if _room.CurrentPlayer.Name != "" && word.Player != _room.CurrentPlayer.Name {
		sendError(w, http.StatusUnauthorized, "Not your turn!")
		return
	}

	s := strings.Split(word.Word, "")

	if len(s) <= 2 {
		sendError(w, http.StatusBadRequest, "Word must begin with "+_room.NextWordStartWith)
		return
	}

	startsWith := strings.Join(s[:2], "")
	if _room.NextWordStartWith != "" && startsWith != _room.NextWordStartWith {
		sendError(w, http.StatusBadRequest, "Word must start with "+_room.NextWordStartWith)
		return
	}

	endsWith := strings.Join(s[len(s)-2:], "")
	_room.Words = append(_room.Words, word)
	_room.NextWordStartWith = endsWith
	_room.CurrentPlayer = serveNextPlayer(&_room)

	fmt.Println(_room)

	updated, err := updateRoom(_room)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "All went to shit while updating room: "+err.Error())
		return
	}

	socket.BroadcastTo("room:"+_room.ID, "room:update", _room)
	sendResp(w, updated)
}

func ClearWords(w http.ResponseWriter, r *http.Request) {
	var roomID = mux.Vars(r)["roomId"]
	room, err := getRoom(roomID)

	if err != nil {
		sendError(w, http.StatusNotFound, "Room not found")
		return
	}

	var _room Room
	json.Unmarshal([]byte(room), &_room)
	_room.Words = []Word{}
	updateRoom(_room)
	sendResp(w, map[string]int{
		"statusCode": http.StatusAccepted,
	})
}
