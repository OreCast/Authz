package main

import (
	"github.com/gin-gonic/gin"
)

type DocsParams struct {
	Login string `json:"login" uri:"login" binding:"required"`
}

// UsersHandler provides access to GET /sites end-point
func UsersHandler(c *gin.Context) {
	var params DocsParams
	if err := c.ShouldBindUri(&params); err == nil {
		user := User{
			LOGIN: params.Login,
		}
		u := getUser(_DB, user)
		c.JSON(200, gin.H{"status": "ok", "user": u})
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
