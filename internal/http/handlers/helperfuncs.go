package handlers

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gary-norman/forum/internal/models"
)

func ErrorMsgs() *models.Errors {
	return models.CreateErrorMessages()
}

func Colors() *models.Colors {
	return models.CreateColors()
}

func IsValidPassword(password string) bool {
	// At least 8 characters
	if len(password) < 8 {
		return false
	}
	// At least one digit
	hasDigit, _ := regexp.MatchString(`[0-9]`, password)
	if !hasDigit {
		return false
	}
	// At least one lowercase letter
	hasLower, _ := regexp.MatchString(`[a-z]`, password)
	if !hasLower {
		return false
	}
	// At least one uppercase letter
	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	return hasUpper
}

func GetTimeSince(created time.Time) string {
	now := time.Now()
	hours := now.Sub(created).Hours()
	var timeSince string
	if hours > 24 {
		timeSince = fmt.Sprintf("%.0f days ago", hours/24)
	} else if hours > 1 {
		timeSince = fmt.Sprintf("%.0f hours ago", hours)
	} else if minutes := now.Sub(created).Minutes(); minutes > 1 {
		timeSince = fmt.Sprintf("%.0f minutes ago", minutes)
	} else {
		timeSince = "just now"
	}
	return timeSince
}

func GetRandomChannel(channelSlice []models.Channel) models.Channel {
	rndInt := rand.IntN(len(channelSlice))
	channel := channelSlice[rndInt]
	return channel
}

func GetRandomUser(userSlice []models.User) models.User {
	rndInt := rand.IntN(len(userSlice))
	user := userSlice[rndInt]
	return user
}

func GetFileName(r *http.Request, fileFieldName, calledBy, imageType string) string {
	// Limit the size of the incoming file to prevent memory issues
	parseErr := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if parseErr != nil {
		log.Printf(ErrorMsgs().Parse, "image", calledBy, parseErr)
		return "noimage"
	}
	// Retrieve the file from form data
	file, handler, retrieveErr := r.FormFile(fileFieldName)
	if retrieveErr != nil {
		log.Printf(ErrorMsgs().RetrieveFile, "image", calledBy, retrieveErr)
		return "noimage"
	}
	defer func(file multipart.File) {
		closeErr := file.Close()
		if closeErr != nil {
			log.Printf(ErrorMsgs().Close, file, calledBy, closeErr)
		}
	}(file)
	// Create a file in the server's local storage
	renamedFile := renameFileWithUUID(handler.Filename)
	fmt.Printf(ErrorMsgs().KeyValuePair, "File Name", renamedFile)
	dst, createErr := os.Create("db/userdata/images/" + imageType + "-images/" + renamedFile)
	if createErr != nil {
		log.Printf(ErrorMsgs().CreateFile, "image", calledBy, createErr)
		return ""
	}
	defer func(dst *os.File) {
		closeErr := dst.Close()
		if closeErr != nil {
			log.Printf(ErrorMsgs().Close, dst, calledBy, closeErr)
		}
	}(dst)
	// Copy the uploaded file data to the server's file
	_, copyErr := io.Copy(dst, file)
	if copyErr != nil {
		fmt.Printf(ErrorMsgs().SaveFile, file, dst, calledBy, copyErr)
		return ""
	}
	return renamedFile
}

func renameFileWithUUID(oldFilePath string) string {
	// Generate a new UUID
	newUUID := models.GenerateToken(16)

	// Split the file name into its base and extension
	base := filepath.Base(oldFilePath)
	ext := filepath.Ext(base)
	// base = base[:len(base)-len(ext)]

	// Create the new file name with the UUID and extension
	newFilePath := filepath.Join(filepath.Dir(oldFilePath), newUUID+ext)

	return newFilePath
}
