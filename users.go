package main

import (
	"fmt"

	"gorm.io/gorm"
)

// User represents user table
type User struct {
	ID        uint
	Login     string
	FirstName string
	LastName  string
	Password  string
	Email     string
	Updated   int64 `gorm:"autoUpdateTime:milli"` // Use unix milli seconds as updating time
	Created   int64 `gorm:"autoCreateTime"`       // Use unix seconds as creating time
}

func getUser(db *gorm.DB, user User) User {
	// Get first matched record
	cond := fmt.Sprintf("FirstName = ? and LastName = ? and Email = ?", user.FirstName, user.LastName, user.Email)
	db.Where(cond).First(&user)
	// SELECT * FROM users WHERE FirstName = ... and LastName = ...
	return user
}

func createUser(db *gorm.DB, user User) (uint, error) {
	result := db.Create(&user) // pass pointer of data to Create
	return user.ID, result.Error
}
