package view

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gary-norman/forum/internal/models"
)

type NotFound interface {
	GetInstance() string
}

func RenderErrorPage[T NotFound](w http.ResponseWriter, data T, error error) {
	var renderedPage bytes.Buffer
	instance := data.GetInstance()
	HTMLstr := models.ToHTMLVar(instance)

	err := Template.ExecuteTemplate(&renderedPage, instance, data)
	if err != nil {
		errorStr := fmt.Sprintf("Error rendering %v: %v", instance, err)
		log.Println(errorStr)
		http.Error(w, errorStr, http.StatusInternalServerError)
		return
	}

	// Send the pre-rendered HTML as JSON
	response := map[string]string{
		HTMLstr: renderedPage.String(),
	}

	w.WriteHeader(http.StatusInternalServerError)
	// Write the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		errorStr := fmt.Sprintf("Error encoding JSON for %v: %v", instance, err)
		http.Error(w, errorStr, http.StatusInternalServerError)
	}

	log.Println(error)
}
