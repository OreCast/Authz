package main

// config module
//
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Configuration stores server configuration parameters
type Configuration struct {
	// web server parts
	Base         string `json:"base"`     // base URL
	LogFile      string `json:"log_file"` // server log file
	Port         int    `json:"port"`     // server port number
	Verbose      int    `json:"verbose"`  // verbose output
	DbUri        string `json:"dburi"`    // database URI
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Domain       string `json:"domain"`

	// OreCast parts
	Cipher string `json:"cipher"` // data-discovery cipher
}

// Config variable represents configuration object
var Config Configuration

// helper function to parse server configuration file
func parseConfig(configFile string) error {
	data, err := os.ReadFile(filepath.Clean(configFile))
	if err != nil {
		log.Println("WARNING: Unable to read", err)
	} else {
		err = json.Unmarshal(data, &Config)
		if err != nil {
			log.Println("ERROR: Unable to parse", err)
			return err
		}
	}

	// default values
	if Config.Port == 0 {
		Config.Port = 8380
	}
	if Config.Cipher == "" {
		Config.Cipher = "aes"
	}
	if Config.DbUri == "" {
		Config.DbUri = "auth.db"
	}
	if Config.ClientId == "" {
		Config.ClientId = "client_id"
	}
	if Config.ClientSecret == "" {
		Config.ClientSecret = "client_secret"
	}
	if Config.Domain == "" {
		Config.Domain = fmt.Sprintf("http://localho:%d", Config.Port)
	}
	return nil
}
