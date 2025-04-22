package view

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gary-norman/forum/internal/models"
)

type PageModel interface {
	GetInstance() string
}

func RenderPageData[T PageModel](w http.ResponseWriter, data T) {
	var renderedPage bytes.Buffer
	instance := data.GetInstance()
	HTMLstr := models.ToHTMLVar(instance)

	err := Template.ExecuteTemplate(&renderedPage, instance, data)
	if err != nil {
		errorStr := fmt.Sprintf("Error rendering %v: %v", instance, err)
		log.Printf(errorStr)
		http.Error(w, errorStr, http.StatusInternalServerError)
		return
	}

	// Send the pre-rendered HTML as JSON
	response := map[string]string{
		HTMLstr: renderedPage.String(),
	}

	// Write the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "error encoding JSON"}`, http.StatusInternalServerError)
	}
}
