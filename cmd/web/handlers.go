package main

import (
	"encoding/json"
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
	emailExists, err := app.users.QueryUserEmailExists(email)
	if emailExists == true {
		er := http.StatusConflict
		http.Error(w, "an account is already registered to that email address", er)
		return
	}
	userExists, _ := app.users.QueryUserNameExists(username)
	if userExists == true {
		er := http.StatusConflict
		http.Error(w, "an account is already registered to that username", er)
		return
	}
	hashedPassword, _ := models.HashPassword(password)
	insertErr := app.users.Insert(
		username,
		email,
		hashedPassword,
		"",
		"",
		"",
		"",
		"")

	if insertErr != nil {
		log.Printf(ErrorMsgs().Register, insertErr)
		http.Error(w, fmt.Sprintf(ErrorMsgs().Register, insertErr), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)

	fprintln, err := fmt.Fprintln(w, "Registration successful")
	if err != nil {
		log.Printf(ErrorMsgs().Printf, err)
		return
	}
	log.Println(fprintln)
}

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	Colors := models.CreateColors()
	login := r.FormValue("username")
	fmt.Printf(Colors.Orange+"Attempting login for "+Colors.White+"%v\n"+Colors.Reset, login)
	fmt.Printf(ErrorMsgs().Divider)
	password := r.FormValue("login_password")
	var user *models.User
	user, getUserErr := app.users.GetUserFromLogin(login, "login")
	if getUserErr != nil {
		log.Printf(ErrorMsgs().NotFound, "either", login, "login > GetUserFromLogin", getUserErr)
	}
	//Notify := models.Notify{
	//	BadPass:      "The passwords do not match.",
	//	RegisterOk:   "Registration successful.",
	//	RegisterFail: "Registration failed.",
	//	BadLogin:     "Username or email address not found",
	//	LoginOk:      "Logged in successfully!",
	//	LoginFail:    "Unable to log in.",
	//}

	if !models.CheckPasswordHash(password, user.HashedPassword) {
		http.Error(w, "incorrect password", http.StatusUnauthorized)
		return
	}
	fmt.Printf(Colors.Green+"Passwords for %v match\n"+Colors.Reset, user.Username)

	// Set Session Token and CSRF Token cookies
	createCookiErr := app.cookies.CreateCookies(w, user)
	if createCookiErr != nil {
		log.Printf(ErrorMsgs().Cookies, "create", createCookiErr)
		http.Error(w, "Failed to create cookies", http.StatusInternalServerError)
		return
	}
	// wait for cookies to be set

	log.Printf(Colors.Green+"Login Successful! Welcome, %v!\n"+Colors.Reset, user.Username)

	//fprintln, err := fmt.Fprintf(w, "Welcome, %v!", user.Username)
	//if err != nil {
	//	log.Print(ErrorMsgs.Login, err)
	//	return
	//}
	//log.Println(fprintln)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *app) logout(w http.ResponseWriter, r *http.Request) {
	Colors := models.CreateColors()
	// Retrieve the cookie
	cookie, cookiErr := r.Cookie("username")
	if cookiErr != nil {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}
	username := cookie.Value
	fmt.Printf(Colors.Orange+"Attempting logout for "+Colors.White+"%v\n"+Colors.Reset, username)
	fmt.Printf(ErrorMsgs().Divider)
	var user *models.User
	user, getUserErr := app.users.GetUserByUsername(username, "logout")
	if getUserErr != nil {
		log.Printf("GetUserByUsername for %v failed with error: %v", username, getUserErr)
	}
	if authErr := app.isAuthenticated(r, username); authErr != nil {
		http.Error(w, authErr.Error(), http.StatusUnauthorized)
		return
	}
	// Delete the Session Token and CSRF Token cookies
	delCookiErr := app.cookies.DeleteCookies(user)
	if delCookiErr != nil {
		log.Printf(ErrorMsgs().Cookies, "delete", delCookiErr)
	}
	fprintln, err := fmt.Fprintln(w, "Logged out successfully!")
	if err != nil {
		return
	}
	log.Println(fprintln)
}

func (app *app) protected(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("username")
	var user *models.User
	user, getUserErr := app.users.GetUserFromLogin(login, "protected")
	if getUserErr != nil {
		log.Printf("protected route for %v failed with error: %v", login, getUserErr)
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if authErr := app.isAuthenticated(r, user.Username); authErr != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fprintf, err := fmt.Fprintf(w, "CSRF Valildation successful! Welcome, %s", user.Username)
	if err != nil {
		log.Print(ErrorMsgs().Protected, user.Username, err)
		return
	}
	log.Println(fprintf)
}

func (app *app) getHome(w http.ResponseWriter, r *http.Request) {
	posts, err := app.posts.All()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	postsWithDaysAgo := make([]models.PostWithDaysAgo, len(posts))

	for index, post := range posts {
		now := time.Now()
		hours := now.Sub(post.Created).Hours()
		var TimeSince string
		if hours > 24 {
			TimeSince = fmt.Sprintf("%.0f days ago", hours/24)
		} else if hours > 1 {
			TimeSince = fmt.Sprintf("%.0f hours ago", hours)
		} else if minutes := now.Sub(post.Created).Minutes(); minutes > 1 {
			TimeSince = fmt.Sprintf("%.0f minutes ago", minutes)
		} else {
			TimeSince = "just now"
		}
		postsWithDaysAgo[index] = models.PostWithDaysAgo{
			Post:      post,
			TimeSince: TimeSince,
		}
	}

	templateData := models.TemplateData{
		Posts:     postsWithDaysAgo,
		Images:    nil,
		Comments:  nil,
		Reactions: nil,
		NotifyPlaceholder: models.NotifyPlaceholder{
			Register: "",
			Login:    "",
		},
	}

	t, err := template.ParseFiles("./assets/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf(ErrorMsgs().Parse, "./assets/templates/index.html", "getHome", err)
		return
	}

	err = t.Execute(w, templateData)
	if err != nil {
		log.Print(ErrorMsgs().Execute, err)
		return
	}
}

func (app *app) createPost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./assets/templates/posts.create.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf(ErrorMsgs().Parse, "./assets/templates/posts.create.html", "createPost", err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf(ErrorMsgs().Execute, err)
		return
	}
}

func (app *app) storePost(w http.ResponseWriter, r *http.Request) {
	user, getUserErr := app.GetLoggedInUser(r)
	if getUserErr != nil {
		http.Error(w, getUserErr.Error(), http.StatusUnauthorized)
		return
	}

	parseErr := r.ParseForm()
	if parseErr != nil {
		http.Error(w, parseErr.Error(), 400)
		log.Printf(ErrorMsgs().Parse, "storePost", parseErr)
		return
	}

	type FormData struct {
		title       string
		content     string
		images      string
		userName    string
		userID      int
		channelName string
		channelID   int
		commentable bool
		isFlagged   bool
	}
	type ChannelData struct {
		ChannelName string `json:"channelName"`
		ChannelID   string `json:"channelID"`
	}
	selectionJSON := r.PostForm.Get("channel")
	if selectionJSON == "" {
		http.Error(w, "No selection provided", http.StatusBadRequest)
		return
	}
	var channelData ChannelData
	if err := json.Unmarshal([]byte(selectionJSON), &channelData); err != nil {
		log.Printf(ErrorMsgs().Unmarshal, selectionJSON, err)
		http.Error(w, "Invalid selection format", http.StatusBadRequest)
		return
	}
	fmt.Printf(ErrorMsgs().KeyValuePair, "channelName", channelData.ChannelName)
	fmt.Printf(ErrorMsgs().KeyValuePair, "channelID", channelData.ChannelID)

	formData := FormData{
		title:       r.PostForm.Get("title"),
		content:     r.PostForm.Get("content"),
		images:      "noimage",
		userName:    user.Username,
		userID:      user.ID,
		channelName: "channelName",
		channelID:   0,
		commentable: false,
		isFlagged:   false,
	}
	if r.PostForm.Get("commentable") != "" {
		formData.commentable = true
	}
	images := r.PostForm.Get("images")
	if images != "" {
		formData.images = images
	}
	formData.channelName = channelData.ChannelName
	formData.channelID, _ = strconv.Atoi(channelData.ChannelID)

	insertErr := app.posts.Insert(
		formData.title,
		formData.content,
		formData.images,
		formData.userName,
		formData.channelName,
		formData.channelID,
		formData.userID,
		formData.commentable,
		formData.isFlagged,
	)

	if insertErr != nil {
		log.Printf(ErrorMsgs().Post, insertErr)
		http.Error(w, insertErr.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
