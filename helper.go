package main

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/nuuls/log"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func sendResp(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		sendResp(w, map[string]interface{}{
			"error":      err.Error(),
			"statusCode": 500,
		})
	}
	w.Write(resp)
}

func sendError(writer http.ResponseWriter, status int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	error := map[string]interface{}{
		"error":      http.StatusText(status),
		"statusCode": status,
	}
	if message != "" {
		error["message"] = message
	}
	sendResp(writer, error)
}

// parseBody parses a request and stores the value in an interface
func parseBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

type middleware func(http.HandlerFunc) http.HandlerFunc

func chainMiddleware(mw ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

func assign(target *map[string]interface{}, ogs ...map[string]interface{}) {
	for _, origin := range ogs {
		for k, v := range origin {
			(*target)[k] = v
		}
	}
}
