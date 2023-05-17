package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	db_key_location = "redis_key.json"
)

var EncryptionKey []byte

type Key struct {
	HostName string
	Addr     string
	Username string
	Password string
	DB       int
	Shell    string
	Key      string
}

func GetKey() Key {
	redis_key_file, err := tryRead()
	if err != nil {
		fmt.Println("Failed to Find Old Redis Key, Please enter new one")
		return newKey()
	}
	return_data := Key{}
	err = json.Unmarshal(redis_key_file, &return_data)
	if err != nil {
		fmt.Println("Failed to Parse Old Key, Overwriting Old One: ", err)
		return newKey()
	}
	HostName = return_data.HostName
	return return_data
}

func newKey() Key {
	file, err := tryCreate()
	if err != nil {
		log.Fatal("Failed To Create File Due to: ", err)
	}
	return_data := Key{}
	return_data.Addr = GetInput("Enter Address of Redis Database, ex: my-redis.cloud.redislabs.com:6379:")
	return_data.DB = GetInt("Enter Database Number(0 is default):")
	return_data.Username = GetInput("Enter User Name Of Database(default is default):")
	return_data.Password = GetInput("Enter Password Of DataBase:")
	return_data.HostName = GetInput("Enter Host Name for YOUR Server:")
	return_data.Shell = strings.Trim(GetInput("Enter What Shell You Want To Use(ex: zsh or bash)"), " ")
	return_data.Key = strings.Trim(GetInput("Enter Your Key"), " ")
	HostName = return_data.HostName
	write_data, err := json.Marshal(return_data)
	if err != nil {
		log.Fatal("FATAL INTERNAL ERROR\nUNABLE TO SET JSON:", err)
	}
	file.Write(write_data)
	err = file.Close()
	if err != nil {
		log.Fatal("Failed to write data due to:", err)
	}
	return return_data
}

func GetInt(message string) int {
	var key string
	for {
		fmt.Println(message)
		fmt.Scan(&key)
		val, err := strconv.Atoi(key)
		if err != nil {
			log.Println("Failed To Parse Int: ", err, "\nPlease Try Again")
		} else {
			return val
		}
	}
}

func GetInput(message string) string {
	var key string
	fmt.Println(message)
	fmt.Scan(&key)
	return key
}

func tryRead() ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		filePath := filepath.Join(homeDir, db_key_location)
		data, err := os.ReadFile(filePath)
		if err == nil {
			return data, nil
		}
	}
	data, err := os.ReadFile(db_key_location)
	if err == nil {
		return data, nil
	}
	filePath := filepath.Join("C:\\", db_key_location)
	data, err = os.ReadFile(filePath)
	if err == nil {
		return data, nil
	}
	return nil, err
}

func tryCreate() (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		filePath := filepath.Join(homeDir, db_key_location)
		data, err := os.Create(filePath)
		if err == nil {
			return data, nil
		}
	}
	data, err := os.Create(db_key_location)
	if err == nil {
		return data, nil
	}
	filePath := filepath.Join("C:\\", db_key_location)
	data, err = os.Create(filePath)
	if err == nil {
		return data, nil
	}
	return nil, err
}
