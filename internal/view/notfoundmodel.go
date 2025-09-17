package view

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gary-norman/forum/internal/models"
)

type ErrorPageData interface {
	GetInstance() string
}

func RenderErrorPage[T ErrorPageData](w http.ResponseWriter, data T, status int, cause error) {
	var renderedPage bytes.Buffer
	instance := data.GetInstance()        // e.g. "user" or "post"
	htmlKey := models.ToHTMLVar(instance) // e.g. "userHTML"

	// Try rendering template
	if err := Template.ExecuteTemplate(&renderedPage, instance, data); err != nil {
		errorStr := fmt.Sprintf("Error rendering %v: %v", instance, err)
		log.Println(errorStr)
		http.Error(w, errorStr, http.StatusInternalServerError)
		return
	}

	// Prepare JSON payload
	response := map[string]any{
		"status": status,
		htmlKey:  renderedPage.String(),
	}

	// Set status code (defaults to 500 if misused)
	if status < 400 {
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)

	// Write JSON to response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		errorStr := fmt.Sprintf("Error encoding JSON for %v: %v", instance, err)
		http.Error(w, errorStr, http.StatusInternalServerError)
		return
	}

	if cause != nil {
		log.Printf("Rendering error page for %v (status %d): %v\n", instance, status, cause)
	}
}
