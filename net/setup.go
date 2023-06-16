package net

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	speedJson "github.com/json-iterator/go"
	"github.com/mattn/go-isatty"
)

const (
	KeyLocation = "redis_key.json"
)

type Key struct {
	HostName string
	Addr     string
	Username string
	Password string
	DB       uint16
	Shell    string
	Key      string
}

func GetKey() Key {
	redisKeyFile, err := tryRead()
	if err != nil {
		fmt.Println("Failed to Find Old Redis Key, Please enter new one")
		return newKey()
	}
	returnData := Key{}
	err = speedJson.ConfigCompatibleWithStandardLibrary.Unmarshal(redisKeyFile, &returnData)
	if err != nil {
		fmt.Println("Failed to Parse Old Key, Overwriting Old One: ", err)
		return newKey()
	}
	return returnData
}

func newKey() Key {
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		log.Fatal("Please set Redis key/create redis key json, exiting(non-tty term)")
	}
	file, err := tryCreate()
	if err != nil {
		log.Fatal("Failed To Create File Due to: ", err)
	}
	returnData := Key{}
	returnData.Addr = GetInput("Enter Address of Redis Database, ex: my-redis.cloud.redislabs.com:6379:")
	returnData.DB = GetInt("Enter Database Number(0 is default):")
	returnData.Username = GetInput("Enter User Name Of Database(default is default):")
	returnData.Password = GetInput("Enter Password Of DataBase:")
	returnData.HostName = GetInput("Enter Host Name for YOUR Server:")
	returnData.Shell = strings.Trim(GetInput("Enter What Shell You Want To Use(ex: zsh or bash)"), " ")
	returnData.Key = strings.Trim(GetInput("Enter Your Key"), " ")
	writeData, err := speedJson.ConfigCompatibleWithStandardLibrary.Marshal(returnData)
	if err != nil {
		log.Fatal("FATAL INTERNAL ERROR\nUNABLE TO SET JSON:", err)
	}
	_, err = file.Write(writeData)
	if err != nil {
		log.Fatal("Failed to Write Data to File d%", err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal("Failed to write data due to:", err)
	}
	return returnData
}

func GetInt(message string) uint16 {
	var key string
	for {
		fmt.Println(message)
		_, err := fmt.Scan(&key)
		if err != nil {
			log.Println("Failed Get int due to ", err)
		}
		val, err := strconv.Atoi(key)
		if err != nil {
			log.Println("Failed To Parse Int: ", err, "\nPlease Try Again")
		} else {
			return uint16(val)
		}
	}
}

func GetInput(message string) string {
	var key string
	fmt.Println(message)
	_, err := fmt.Scan(&key)
	if err != nil {
		log.Fatal("Failed to Get input due to: ", err)
	}
	return key
}

func tryRead() ([]byte, error) {
	data, err := os.ReadFile(KeyLocation)
	if err == nil {
		return data, nil
	}
	homeDir, err := os.UserHomeDir()
	if err == nil {
		filePath := filepath.Join(homeDir, KeyLocation)
		data, err := os.ReadFile(filePath)
		if err == nil {
			return data, nil
		}
	}
	filePath := filepath.Join("C:\\\\", KeyLocation)
	data, err = os.ReadFile(filePath)
	if err == nil {
		return data, nil
	}
	return nil, err
}

func tryCreate() (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		filePath := filepath.Join(homeDir, KeyLocation)
		data, err := os.Create(filePath)
		if err == nil {
			return data, nil
		}
	}
	data, err := os.Create(KeyLocation)
	if err == nil {
		return data, nil
	}
	filePath := filepath.Join("C:\\", KeyLocation)
	data, err = os.Create(filePath)
	if err == nil {
		return data, nil
	}
	return nil, err
}
