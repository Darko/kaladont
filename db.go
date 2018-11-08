package main

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func initRedis() *redis.Pool {
	p := &redis.Pool{}

	p.Dial = func() (redis.Conn, error) {
		fmt.Println("Connecting to redis")
		conn, err := redis.Dial("tcp", "localhost:6379")
		if err == nil {
			return conn, nil
		}
		return nil, errors.New("failed to connect to redis")
	}

	return p
}

// Db xD
type Db struct {
	Conn redis.Conn
}

// Get key
func (db *Db) Get(key string) (data interface{}, err error) {
	return db.Conn.Do("GET", key)
	// return data, err
}

// Set key
func (db *Db) Set(key string, d interface{}) (data interface{}, err error) {
	return db.Conn.Do("SET", key, d)
}

// Delete key
func (db *Db) Delete(key string) (data interface{}, err error) {
	return db.Conn.Do("DEL", key)
}
