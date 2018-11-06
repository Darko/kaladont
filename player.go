package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Player struct
type Player struct {
	Name  string `json:"name"`
	Score int32  `json:"score"`
}

func findPlayer(p []Player, name string) (interface{}, error) {
	for _, item := range p {
		if item.Name == name {
			return item, nil
		}
	}

	return nil, errors.New("Player not found")
}

// JoinRoom controller
func JoinRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var roomID = vars["roomId"]
	var playerName = vars["name"]

	if playerName == "" {
		sendError(w, 400, "Missing parameter: name")
		return
	}

	result, err := getRoom(roomID)

	if err != nil {
		sendError(w, 404, "Room not found")
		return
	}

	var room Room
	json.Unmarshal([]byte(result), &room)

	_, err = findPlayer(room.Players, playerName)

	if err != nil {
		player := Player{playerName, 0}
		room.Players = append(room.Players, player)

		updated, _ := updateRoom(room)
		sendResp(w, updated)
		return
	}

	sendError(w, 400, "Username is taken")
}

func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	var authHeader = r.Header.Get("Authorization")
	var apiKey = strings.Split(authHeader, " ")[1]

	fmt.Println(strings.Fields(authHeader))

	if apiKey == "" {
		sendError(w, 401, "")
	}
}
