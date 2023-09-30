package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// examples: https://go.dev/doc/tutorial/web-service-gin

var _DB *gorm.DB

// helper function to setup our server router
func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// GET routes
	r.GET("/user/:login", UsersHandler)

	// POST routes
	r.POST("/user", UsersPostHandler)

	return r
}

func Server(configFile string) {
	db, err := initDB("sqlite")
	if err != nil {
		log.Fatal(err)
	}
	_DB = db
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := setupRouter()
	sport := fmt.Sprintf(":%d", Config.Port)
	log.Printf("Start HTTP server %s", sport)
	r.Run(sport)
}
