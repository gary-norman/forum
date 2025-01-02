package main

import (
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func isValidPassword(password string) bool {
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
	if !hasUpper {
		return false
	}

	return true
}

func (app *app) register(w http.ResponseWriter, r *http.Request) {
	ErrorMsgs := models.CreateErrorMessages()
	username := r.FormValue("register_user")
	email := r.FormValue("register_email")
	validEmail, _ := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, email)
	password := r.FormValue("register_password")
	if len(username) < 5 || len(username) > 16 {
		er := http.StatusNotAcceptable
		http.Error(w, "username must be between 5 and 16 characters", er)
	}
	if isValidPassword(password) != true {
		er := http.StatusNotAcceptable
		http.Error(w, "password must contain at least one number and one uppercase and lowercase letter,"+
			"and at least 8 or more characters", er)
	}
	if validEmail != true {
		er := http.StatusNotAcceptable
		http.Error(w, "please enter a valid email address", er)
	}
	exists := app.users.QueryUserEmailExists(email)
	if exists == true {
		er := http.StatusConflict
		http.Error(w, "an account is already registered to that email address", er)
		return
	}

	if app.users.QueryUserNameExists(username) {
		er := http.StatusConflict
		http.Error(w, "username already exists", er)
		return
	}
	hashedPassword, _ := models.HashPassword(password)
	err := app.users.Insert(
		username,
		email,
		hashedPassword,
		"",
		"",
		"",
		"",
		"")

	if err != nil {
		log.Printf(ErrorMsgs.Register, err)
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)

	fprintln, err := fmt.Fprintln(w, "Registration successful")
	if err != nil {
		log.Print(ErrorMsgs.Register, err)
		return
	}
	log.Println(fprintln)
}

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("login_username")
	email := r.FormValue("login_username")
	password := r.FormValue("login_password")
	if app.users.QueryUserNameExists(username) == true {
		user, _ := app.users.GetUserByUsername(username)
		if !models.CheckPasswordHash(password, user.HashedPassword) {
			http.Error(w, "incorrect password", http.StatusUnauthorized)
			return
		}
	}

	if app.users.QueryUserEmailExists(email) == true {
		user, _ := app.users.GetUserByEmail(email)
		if !models.CheckPasswordHash(password, user.HashedPassword) {
			http.Error(w, "incorrect password", http.StatusUnauthorized)
			return
		}
	}

	if app.users.QueryUserNameExists(username) == false && app.users.QueryUserEmailExists(email) == false {
		http.Error(w, "username or email not found", http.StatusUnauthorized)
		return
	}

	sessionToken := models.GenerateToken(32)
	csrfToken := models.GenerateToken(32)

	// Set Session Token and CSRF Token cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: false,
	})

	// Store tokens in the database
	err := app.users.Insert(
		"",
		"",
		"",
		sessionToken,
		csrfToken,
		"",
		"",
		"")

	ErrorMsgs := models.CreateErrorMessages()

	if err != nil {
		log.Printf(ErrorMsgs.Login, err)
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)

	fprintln, err := fmt.Fprintln(w, "Logged in successfully!")
	if err != nil {
		log.Print(ErrorMsgs.Login, err)
		return
	}
	log.Println(fprintln)
}

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
		var TimeSince string
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
		log.Printf(ErrorMsgs.Parse, "storePost", err)
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
