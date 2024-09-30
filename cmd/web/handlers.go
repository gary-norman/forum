package main

import (
	"github.com/gary-norman/forum/internal/models"
	"html/template"
	"log"
	"net/http"
)

func (app *app) getHome(w http.ResponseWriter, r *http.Request) {
	ErrorMsgs := models.CreateErrorMessages()
	posts, err := app.posts.All()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t, err := template.ParseFiles("./assets/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Print(ErrorMsgs.Parse, "./assets/templates/index.html", "template.ParseFiles", err)
		return
	}

	err = t.Execute(w, map[string]any{"Posts": posts})
	if err != nil {
		log.Print(ErrorMsgs.Execute, err)
		return
	}
}
