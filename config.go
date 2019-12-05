package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// BirthDayBotConfig struct to hold the configuration
type BirthDayBotConfig struct {
	MMIncomingWebhook    string
	GoogleCalendarID     string
	GoogleCalendarAPIKey string
}

// Config holds the configuration
var Config = &BirthDayBotConfig{}

func loadConfig(fileName string) error {
	fileName = findConfigFile(fileName)
	log.Println("Loading " + fileName)

	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Error opening config file=" + fileName + ", err=" + err.Error())
		return err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(Config)
	if err != nil {
		log.Println("Error decoding config file=" + fileName + ", err=" + err.Error())
		return err
	}

	return nil
}

func findConfigFile(fileName string) string {
	if _, err := os.Stat("/tmp/" + fileName); err == nil {
		fileName, _ = filepath.Abs("/tmp/" + fileName)
	} else if _, err := os.Stat("./config/" + fileName); err == nil {
		fileName, _ = filepath.Abs("./config/" + fileName)
	} else if _, err := os.Stat("../config/" + fileName); err == nil {
		fileName, _ = filepath.Abs("../config/" + fileName)
	} else if _, err := os.Stat(fileName); err == nil {
		fileName, _ = filepath.Abs(fileName)
	}

	return fileName
}
