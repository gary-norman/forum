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
	htmlkey := models.ToHTMLVar(instance)

	if err := Template.ExecuteTemplate(&renderedPage, instance, data); err != nil {
		errorStr := fmt.Sprintf("Error rendering %v: %v", instance, err)
		log.Println(errorStr)
		http.Error(w, errorStr, http.StatusInternalServerError)
		return
	}

	// Send the pre-rendered HTML as JSON
	response := map[string]string{
		htmlkey: renderedPage.String(),
	}

	// Write the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		RenderErrorPage(w, models.NotFoundLocation("page"), 500, models.EncodeError(fmt.Sprintf("response for %v", instance), "RenderPageData", err))
	}
}
