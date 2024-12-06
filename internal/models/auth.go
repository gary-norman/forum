package models

import (
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

// Users Key is the username
var Users = map[string]Login{}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
