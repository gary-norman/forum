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
		log.Printf(ErrorMsgs.Parse, "./assets/templates/index.html", "getHome", err)
		return
	}

	err = t.Execute(w, map[string]any{"Posts": posts})
	if err != nil {
		log.Print(ErrorMsgs.Execute, err)
		return
	}
}

func (app *app) createPost(w http.ResponseWriter, r *http.Request) {
	ErrorMsgs := models.CreateErrorMessages()
	t, err := template.ParseFiles("./assets/templates/posts.create.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf(ErrorMsgs.Parse, "./assets/templates/posts.create.html", "createPost", err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf(ErrorMsgs.Execute, err)
		return
	}
}

func (app *app) storePost(w http.ResponseWriter, r *http.Request) {
	ErrorMsgs := models.CreateErrorMessages()
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Printf(ErrorMsgs.Parse, "./assets/templates/posts.create.html", "storePost", err)
		return
	}
	err = app.posts.Insert(
		r.PostForm.Get("title"),
		r.PostForm.Get("content"),
	)
	if err != nil {
		log.Printf(ErrorMsgs.Post, err)
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
