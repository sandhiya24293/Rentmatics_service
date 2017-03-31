package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port       string `json:"port"`
	Https_port string `json:"https_port"`

	SSL  bool   `json:"ssl"`
	Key  string `json:"key"`
	Cert string `json:"cert"`
}

var (
	config = Config{}
)

func getConfig() bool {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error during opening configuration file: ", err)
		return false
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error during decoding configuration file: ", err)
		return false
	}
	return true
}
