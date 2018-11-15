package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func findPlayer(p []Player, name string) (Player, int, error) {
	for i, item := range p {
		if item.Name == name {
			return item, i, nil
		}
	}

	return Player{}, -1, errors.New("Player not found")
}

func removePlayerByName(p []Player, name string) []Player {
	_, ind, _ := findPlayer(p, name)

	if ind > -1 {
		p = append(p[:ind], p[ind+1:]...)
	}

	return p
}

func addScore(players *[]Player, playerName string) {
	_, ind, _ := findPlayer((*players), playerName)
	(*players)[ind].Score += 3
}

// JoinRoom controller
func JoinRoom(w http.ResponseWriter, r *http.Request) {
	var roomID = mux.Vars(r)["roomId"]
	var player Player
	var err = parseBody(r, &player)

	if err != nil {
		sendError(w, 500, err.Error())
		return
	}

	if player.Name == "" {
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

	_, _, err = findPlayer(room.Players, player.Name)

	if err != nil {
		player := Player{randomID(16), player.Name, 0}
		room.Players = append(room.Players, player)

		updated, _ := updateRoom(room)
		go socket.BroadcastTo("room:"+room.ID, "player:join", map[string]interface{}{
			"player": player,
		})
		sendResp(w, updated)
		return
	}

	sendError(w, 400, "Username is taken")
}

// LeaveRoom controller
func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	var roomID = mux.Vars(r)["roomId"]
	p := map[string]string{}
	err := parseBody(r, &p)

	if err != nil {
		sendError(w, 500, err.Error())
		return
	}

	room, err := getRoom(roomID)
	if err != nil {
		println(err.Error())
	}

	var _room Room
	fmt.Println(room)
	err = json.Unmarshal([]byte(room), &_room)

	if err != nil {
		println(err.Error())
	}

	_room.Players = removePlayerByName(_room.Players, p["name"])
	_, err = updateRoom(_room)

	if err != nil {
		sendError(w, 500, err.Error())
		return
	}

	sendResp(w, map[string]int{
		"statusCode": 201,
	})
}
