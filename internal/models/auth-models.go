package models

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Email          string
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

type Session struct {
	Username string
	Expires  time.Time
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
