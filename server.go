package main

import (
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

// examples: https://go.dev/doc/tutorial/web-service-gin

// _DB defines gorm DB pointer
var _DB *gorm.DB

var _oauthServer *server.Server

// helper function to setup our server router
func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// GET routes
	r.GET("/user/:login", UsersHandler)
	r.GET("/oauth/token", TokenHandler)

	// POST routes
	r.POST("/user", RegistryUserHandler)
	r.POST("/oauth/authorize", AuthzHandler)

	return r
}

func Server(configFile string) {
	db, err := initDB("sqlite")
	if err != nil {
		log.Fatal(err)
	}
	_DB = db
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// setup oauth parts
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	manager.MapAccessGenerate(
		generates.NewJWTAccessGenerate(
			"", []byte(Config.ClientId), jwt.SigningMethodHS512))
	//     manager.MapAccessGenerate(generates.NewAccessGenerate())

	clientStore := store.NewClientStore()
	clientStore.Set(Config.ClientId, &models.Client{
		ID:     Config.ClientId,     // The client id being passed in
		Secret: Config.ClientSecret, // The client secret being passed in
		Domain: Config.Domain,       // The domain of the redirect url
	})
	manager.MapClientStorage(clientStore)
	_oauthServer = server.NewServer(server.NewConfig(), manager)
	_oauthServer.SetAllowGetAccessRequest(true)
	_oauthServer.SetClientInfoHandler(server.ClientFormHandler)

	r := setupRouter()
	sport := fmt.Sprintf(":%d", Config.Port)
	log.Printf("Start HTTP server %s", sport)
	r.Run(sport)
}
