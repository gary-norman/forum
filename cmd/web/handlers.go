package main

import (
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (app *app) register(w http.ResponseWriter, r *http.Request) {
	ErrorMsgs := models.CreateErrorMessages()
	username := r.FormValue("username")
	//email := r.FormValue("email")
	password := r.FormValue("password")
	if len(username) < 5 || len(password) < 8 {
		er := http.StatusNotAcceptable
		http.Error(w, "invalid username or password", er)
	}
	if _, ok := models.Users[username]; ok {
		er := http.StatusConflict
		http.Error(w, "username already exists", er)
		return
	}
	hashedPassword, _ := models.HashPassword(password)
	models.Users[username] = models.Login{
		HashedPassword: hashedPassword,
	}
	fprintln, err := fmt.Fprintln(w, "Registration successful")
	if err != nil {
		log.Print(ErrorMsgs.Register, err)
		return
	}
	log.Println(fprintln)
}

func (app *app) login(w http.ResponseWriter, r *http.Request) {}

func (app *app) logout(w http.ResponseWriter, r *http.Request) {}

func (app *app) protected(w http.ResponseWriter, r *http.Request) {}

func (app *app) getHome(w http.ResponseWriter, r *http.Request) {
	ErrorMsgs := models.CreateErrorMessages()
	posts, err := app.posts.All()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	postsWithDaysAgo := make([]models.PostWithDaysAgo, len(posts))

	for index, post := range posts {
		fmt.Printf("index: %v, post: %v\n", index, post.Created)
		now := time.Now()
		hours := now.Sub(post.Created).Hours()
		TimeSince := ""
		if hours > 24 {
			TimeSince = fmt.Sprintf("%.0f days ago.", hours/24)
			fmt.Printf("Hours: %v, days: %v\n", hours, hours/24)
			fmt.Printf("Timesince: %v\n", post.TimeSince)
		} else if hours > 1 {
			TimeSince = fmt.Sprintf("%.0f hours ago.", hours)
			fmt.Printf("Hours: %v\n", hours)
			fmt.Printf("Timesince: %v\n", post.TimeSince)
		} else if minutes := now.Sub(post.Created).Minutes(); minutes > 1 {
			TimeSince = fmt.Sprintf("%.0f minutes ago.", minutes)
			fmt.Printf("Minutes: %v\n", minutes)
			fmt.Printf("Timesince: %v\n", post.TimeSince)
		} else {
			TimeSince = "just now"
			fmt.Printf("Timesince: %v\n", post.TimeSince)
		}
		postsWithDaysAgo[index] = models.PostWithDaysAgo{
			Post:      post,
			TimeSince: TimeSince,
		}
	}

	data := struct {
		Posts []models.PostWithDaysAgo
	}{
		Posts: postsWithDaysAgo,
	}

	t, err := template.ParseFiles("./assets/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf(ErrorMsgs.Parse, "./assets/templates/index.html", "getHome", err)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Print(ErrorMsgs.Execute, err)
		return
	}
	for _, post := range posts {
		fmt.Printf("Timesince 2: %v\n", post.TimeSince)
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

	// Get the 'channel' value as a string
	channelStr := r.PostForm.Get("channel")
	// Convert the string to an integer
	channel, err := strconv.Atoi(channelStr)
	if err != nil {
		http.Error(w, "You must be a member of this channel to do that.", http.StatusBadRequest)
		return
	}

	// Get the 'author' value as a string
	authorStr := r.PostForm.Get("author")
	// Convert the string to an integer
	author, err := strconv.Atoi(authorStr)
	if err != nil {
		http.Error(w, "You must be logged in to do that.", http.StatusBadRequest)
		return
	}

	type FormData struct {
		commentable bool
		images      string
	}
	formData := FormData{
		commentable: false,
		images:      "noimage",
	}
	if r.PostForm.Get("commentable") != "" {
		formData.commentable = true
	}
	images := r.PostForm.Get("images")
	if images != "" {
		formData.images = images
	}

	err = app.posts.Insert(
		r.PostForm.Get("title"),
		r.PostForm.Get("content"),
		formData.images,
		channel,
		author,
		formData.commentable,
	)

	if err != nil {
		log.Printf(ErrorMsgs.Post, err)
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
