package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DocsParams struct {
	Login string `json:"login" uri:"login" binding:"required"`
}
type UserParams struct {
	Login    string `json:"login" uri:"login" binding:"required"`
	Password string `json:"password" uri:"password" binding:"required"`
}

// UsersHandler provides access to GET /sites end-point
func UsersHandler(c *gin.Context) {
	var params DocsParams
	if err := c.ShouldBindUri(&params); err == nil {
		user := User{
			LOGIN: params.Login,
		}
		if u, err := getUser(_DB, user); err == nil {
			c.JSON(200, gin.H{"status": "ok", "user": u})
		} else {
			c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
		}
	} else {
		c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
	}
}

// UserRequest represents user form request
type UserRequest struct {
	Login     string `json:"login"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

// UsersPostHandler provides access to POST /sites end-point
func UsersPostHandler(c *gin.Context) {
	var form UserRequest
	err := c.BindJSON(&form)
	if err == nil {
		// create new user in DB
		user := User{
			LOGIN:      form.Login,
			FIRST_NAME: form.FirstName,
			LAST_NAME:  form.LastName,
			PASSWORD:   form.Password,
			EMAIL:      form.Email,
		}
		uid, err := createUser(_DB, user)
		if err != nil {
			c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
		} else {
			c.JSON(200, gin.H{"status": "ok", "uid": uid})
		}
	} else {
		c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
	}
}

// TokenHandler provides access to GET /oauth/token end-point
func TokenHandler(c *gin.Context) {
	err := _oauthServer.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		log.Println("ERROR: oauth server error", err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}
}

// AuthzHandler provides access to POST /oauth/authorize end-point
func AuthzHandler(c *gin.Context) {
	//     err := _oauthServer.HandleAuthorizeRequest(c.Writer, c.Request)
	//     if err != nil {
	//         log.Println("ERROR: oauth server error", err)
	//         http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	//     }
	var params UserParams
	if err := c.BindJSON(&params); err == nil {
		user := User{
			LOGIN:    params.Login,
			PASSWORD: params.Password,
		}
		if u, err := getUser(_DB, user); err == nil {
			err = _oauthServer.HandleAuthorizeRequest(c.Writer, c.Request)
			if err != nil {
				log.Println("ERROR: oauth server error", err)
				http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			}
			log.Println("INFO: found user", u)
			authServerURL := "http://localhost:8380"
			rurl := fmt.Sprintf("%s/oauth/token?client_id=client_id&client_secret=client_secret&grant_type=client_credentials&scope=read", authServerURL)
			//             rurl := fmt.Sprintf("/oauth/token?client_id=client_id&client_secret=client_secret&grant_type=client_credentials&scope=read")
			r, e := http.NewRequest("GET", rurl, nil)
			r.Header.Add("Accept", "application/json")
			r.Header.Add("Content-type", "application/json")
			if e == nil {
				client := http.Client{}
				resp, err := client.Do(r)
				if err != nil {
					log.Println("ERROR: unable to make new HTTP request", rurl, err)
				} else {
					log.Println("INFO", resp.Status, resp.StatusCode)
					defer resp.Body.Close()
					data, err := io.ReadAll(resp.Body)
					log.Println("DATA", string(data), err)
				}
				//                 http.Redirect(c.Writer, r, rurl, http.StatusFound)
			} else {
				log.Println("ERROR: unable to make new HTTP request", r)
			}
			//             config := oauth2.Config{
			//                 ClientID:     "client_id",
			//                 ClientSecret: "client_secret",
			//                 Scopes:       []string{"all"},
			//                 RedirectURL:  "http://localhost:8380/oauth",
			//                 Endpoint: oauth2.Endpoint{
			//                     AuthURL:  authServerURL + "/oauth/authorize",
			//                     TokenURL: authServerURL + "/oauth/token",
			//                 },
			//             }
			//             if token, err := config.PasswordCredentialsToken(context.Background(), u.LOGIN, u.PASSWORD); err == nil {
			//                 c.JSON(200, gin.H{"status": "ok", "token": token})
			//             } else {
			//                 c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
			//             }
		} else {
			c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
		}
	} else {
		c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
	}
}
