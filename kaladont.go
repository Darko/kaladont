package main

import (
	"errors"

	"github.com/garyburd/redigo/redis"
)

func initRedis() *redis.Pool {
	p := &redis.Pool{}

	p.Dial = func() (redis.Conn, error) {
		conn, err := redis.Dial("tcp", "localhost:6379")
		if err == nil {
			return conn, nil
		}
		return nil, errors.New("failed to connect to redis")
	}

	return p
}

// Kaladont struct
type Kaladont struct {
	redis *redis.Pool
}

func initKaladont() *Kaladont {
	kt := &Kaladont{
		redis: initRedis(),
	}
	return kt
}
