package redis

import (
	"fmt"
	"log"
	"os"
)

const (
	db_key_location = "redis_key"
)

type Key interface {
	Addr string,
	Password string,
	DB int
}

func get_key() string {
	redis_key_file, err := os.ReadFile(db_key_location)
	if err != nil {
		log.Print("Failed to Find Old Redis Key, Please enter new one: ")
		file, err := os.Create(db_key_location)
		if err != nil {
			log.Fatal("Failed To Create File Due to: ", err)
		}
		var key string
		fmt.Scan(&key)
		file.Write([]byte(key))
		err = file.Close()
		if err != nil {
			log.Fatal("Failed to write data due to:", err)
		}
		return key
	}
	return string(redis_key_file)
}
