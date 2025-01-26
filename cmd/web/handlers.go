package main

import (
	"encoding/json"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
		"noimage",
		"default.png",
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
	if username == "" {
		log.Printf(ErrorMsgs().KeyValuePair, "aborting logout:", "no user is logged in")
		return
	}
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
	delCookiErr := app.cookies.DeleteCookies(w, user)
	if delCookiErr != nil {
		log.Printf(ErrorMsgs().Cookies, "delete", delCookiErr)
	}
	log.Println("Check1 /")
	log.Println(Colors.Green + "Logged out successfully!")

	//http.Redirect(w, r, "/", http.StatusFound)

	//w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusNoContent)
	//app.getHome(w, r)
	log.Println("Check2 /")
	//http.Redirect(w, r, "/", http.StatusFound)
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
	userLoggedIn := true
	posts, postsErr := app.posts.All()
	if postsErr != nil {
		http.Error(w, postsErr.Error(), 500)
		return
	}
	// Retrieve total likes and dislikes for each post
	for i, post := range posts {
		likes, dislikes, likesErr := app.reactions.CountReactions(post.ChannelID, post.ID, 0) // Pass 0 for CommentID if it's a post
		fmt.Printf("PostID: %v, Likes: %v, Disikes: %v\n", posts[i].ID, likes, dislikes)
		if likesErr != nil {
			log.Printf("Error counting reactions: %v", likesErr)
			likes, dislikes = 0, 0 // Default values if there is an error
		}
		posts[i].Likes = likes
		posts[i].Dislikes = dislikes
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

	currentUser, currentUserErr := app.GetLoggedInUser(r)
	if currentUserErr != nil {
		log.Printf(ErrorMsgs().NotFound, "user", "current user", "GetLoggedInUser", currentUserErr)
		log.Printf(ErrorMsgs().KeyValuePair, "Current user", currentUser)
		userLoggedIn = false
	}
	var ownedChannels []models.Channel
	var channelsErr error

	currentUserName := "nouser"
	var currentUserAvatar string
	var currentUserBio string
	if userLoggedIn == true {
		currentUserName = currentUser.Username
		currentUserAvatar = currentUser.Avatar
		currentUserBio = currentUser.Description
		ownedChannels, channelsErr = app.channels.OwnedByCurrentUser(currentUser.ID)
		if channelsErr != nil {
			log.Printf(ErrorMsgs().Query, "user channels", channelsErr)
		}
	}

	//fmt.Printf(ErrorMsgs().KeyValuePair, "Owned Channels", ownedChannels)
	//fmt.Printf(ErrorMsgs().KeyValuePair, "Current user", currentUser)
	//fmt.Printf(ErrorMsgs().KeyValuePair, "currentUserAvatar", currentUserAvatar)

	templateData := models.TemplateData{
		OwnedChannels:   ownedChannels,
		CurrentUser:     currentUser,
		CurrentUserName: currentUserName,
		Posts:           postsWithDaysAgo,
		Avatar:          currentUserAvatar,
		Bio:             currentUserBio,
		Images:          nil,
		Comments:        nil,
		Reactions:       nil,
		NotifyPlaceholder: models.NotifyPlaceholder{
			Register: "",
			Login:    "",
		},
	}
	models.JsonError(templateData)

	tpl, err := GetTemplate()
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf(ErrorMsgs().Parse, "./assets/templates/index.html", "getHome", err)
		return
	}

	t := tpl.Lookup("index.html")
	if t == nil {
		http.Error(w, "template not found: index.html", 500)
		log.Printf("Template not found: index.html")
		return
	}

	data := templateData
	execErr := t.Execute(w, data)
	if execErr != nil {
		log.Printf(ErrorMsgs().Execute, execErr)
		return
	}
}

func (app *app) editUserDetails(w http.ResponseWriter, r *http.Request) {
	user, getUserErr := app.GetLoggedInUser(r)
	if getUserErr != nil {
		log.Printf(ErrorMsgs().NotFound, "user", "current user", "GetLoggedInUser", getUserErr)
		return
	}
	if user.Username == "" {
		log.Printf(ErrorMsgs().NotFound, "user", "current user", "GetLoggedInUser", getUserErr)
		return
	}
	parseErr := r.ParseMultipartForm(10 << 20)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), 400)
		log.Printf(ErrorMsgs().Parse, "editUserDetails", parseErr)
		return
	}
	user.Avatar = GetFileName(r, "file-drop", "editUserDetails", "user")
	description := r.FormValue("bio")
	if description != "" {
		user.Description = description
	}
	editErr := app.users.Edit(user)
	if editErr != nil {
		log.Printf(ErrorMsgs().Edit, user.Username, "EditUserDetails", editErr)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *app) createPost(w http.ResponseWriter, r *http.Request) {
	t, parseErr := template.ParseFiles("./assets/templates/posts.create.html")
	if parseErr != nil {
		http.Error(w, parseErr.Error(), 500)
		log.Printf(ErrorMsgs().Parse, "./assets/templates/posts.create.html", "createPost", parseErr)
		return
	}

	execErr := t.Execute(w, nil)
	if execErr != nil {
		log.Printf(ErrorMsgs().Execute, execErr)
		return
	}
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
	base = base[:len(base)-len(ext)]

	// Create the new file name with the UUID and extension
	newFilePath := filepath.Join(filepath.Dir(oldFilePath), newUUID+ext)

	return newFilePath
}

func (app *app) storePost(w http.ResponseWriter, r *http.Request) {
	user, getUserErr := app.GetLoggedInUser(r)
	if getUserErr != nil {
		http.Error(w, getUserErr.Error(), http.StatusUnauthorized)
		return
	}
	parseErr := r.ParseMultipartForm(10 << 20)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), 400)
		log.Printf(ErrorMsgs().Parse, "storePost", parseErr)
		return
	}
	// get channel name
	selectionJSON := r.PostForm.Get("channel")
	if selectionJSON == "" {
		http.Error(w, "No selection provided", http.StatusBadRequest)
		return
	}
	var channelData models.ChannelData
	if err := json.Unmarshal([]byte(selectionJSON), &channelData); err != nil {
		log.Printf(ErrorMsgs().Unmarshal, selectionJSON, err)
		http.Error(w, "Invalid selection format", http.StatusBadRequest)
		return
	}
	fmt.Printf(ErrorMsgs().KeyValuePair, "channelName", channelData.ChannelName)
	fmt.Printf(ErrorMsgs().KeyValuePair, "channelID", channelData.ChannelID)
	fmt.Printf(ErrorMsgs().KeyValuePair, "commentable", r.PostForm.Get("commentable"))

	createPostData := models.Post{
		Title:         r.PostForm.Get("title"),
		Content:       r.PostForm.Get("content"),
		Images:        "",
		Author:        user.Username,
		AuthorID:      user.ID,
		AuthorAvatar:  user.Avatar,
		ChannelName:   "channelName",
		ChannelID:     0,
		IsCommentable: false,
		IsFlagged:     false,
	}
	fmt.Printf(ErrorMsgs().KeyValuePair, "authorAvatar", createPostData.AuthorAvatar)
	if r.PostForm.Get("commentable") == "on" {
		createPostData.IsCommentable = true
	}
	createPostData.Images = GetFileName(r, "file-drop", "storePost", "post")
	createPostData.ChannelName = channelData.ChannelName
	createPostData.ChannelID, _ = strconv.Atoi(channelData.ChannelID)

	insertErr := app.posts.Insert(
		createPostData.Title,
		createPostData.Content,
		createPostData.Images,
		createPostData.Author,
		createPostData.ChannelName,
		createPostData.AuthorAvatar,
		createPostData.ChannelID,
		createPostData.AuthorID,
		createPostData.IsCommentable,
		createPostData.IsFlagged,
	)

	if insertErr != nil {
		log.Printf(ErrorMsgs().Post, insertErr)
		http.Error(w, insertErr.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *app) storeChannel(w http.ResponseWriter, r *http.Request) {
	user, getUserErr := app.GetLoggedInUser(r)
	if getUserErr != nil {
		http.Error(w, getUserErr.Error(), http.StatusUnauthorized)
		return
	}
	parseErr := r.ParseMultipartForm(10 << 20)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), 400)
		log.Printf(ErrorMsgs().Parse, "storeChannel", parseErr)
		return
	}

	createChannelData := models.Channel{
		OwnerID:     user.ID,
		Name:        r.PostForm.Get("name"),
		Description: r.PostForm.Get("description"),
		Avatar:      "noimage",
		Banner:      "default.png",
		Privacy:     false,
		IsFLagged:   false,
		IsMuted:     false,
	}
	if r.PostForm.Get("privacy") == "on" {
		createChannelData.Privacy = true
	}
	createChannelData.Avatar = GetFileName(r, "file-drop", "storeChannel", "channel")

	insertErr := app.channels.Insert(
		createChannelData.OwnerID,
		createChannelData.Name,
		createChannelData.Description,
		createChannelData.Avatar,
		createChannelData.Banner,
		createChannelData.Privacy,
		createChannelData.IsFLagged,
		createChannelData.IsMuted,
	)

	if insertErr != nil {
		log.Printf(ErrorMsgs().Post, insertErr)
		http.Error(w, insertErr.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *app) storeReaction(w http.ResponseWriter, r *http.Request) {

	//log.Printf("using storeReaction()")

	// Check if the method is POST, otherwise return Method Not Allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Variable to hold the decoded data
	var reactionData models.Reaction

	// Decode the JSON request body into the reactionData variable
	err := json.NewDecoder(r.Body).Decode(&reactionData)
	if err != nil {
		// Use the error message from the Errors struct for decoding errors
		http.Error(w, fmt.Sprintf(ErrorMsgs().Parse, "storeReaction", err), http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	//// Validate that at least one of reactedPostID or reactedCommentID is non-zero
	if (reactionData.ReactedPostID == nil || *reactionData.ReactedPostID == 0) && (reactionData.ReactedCommentID == nil || *reactionData.ReactedCommentID == 0) {
		http.Error(w, "You must react to either a post or a comment", http.StatusBadRequest)
		return
	}

	postID, commentID := 0, 0

	if reactionData.ReactedPostID != nil {
		//log.Println("ReactedPostID:", *reactionData.ReactedPostID)
		postID = *reactionData.ReactedPostID
	}

	if reactionData.ReactedCommentID != nil {
		//log.Printf("ReactedCommentID: %d", *reactionData.ReactedPostID)
		commentID = *reactionData.ReactedCommentID
	}

	// Check if the user already reacted (like/dislike) and update or delete the reaction if needed
	existingReaction, err := app.reactions.CheckExistingReaction(reactionData.AuthorID, postID, commentID)
	if err != nil {
		// Use your custom error message for fetching errors
		http.Error(w, fmt.Sprintf(ErrorMsgs().Read, "storeReaction", err), http.StatusInternalServerError)
		log.Printf("Error fetching existing reaction: %v", err)
		return
	}

	//if existingReaction != nil {
	//	//log.Printf("Existing Reaction: %+v", existingReaction)
	//}

	// If there is an existing reaction, toggle it (i.e., remove it if the user reacts again to the same thing)
	if existingReaction != nil {
		// If the user likes a post or comment again, remove the like/dislike (toggle behavior)
		if existingReaction.Liked == reactionData.Liked && existingReaction.Disliked == reactionData.Disliked {
			// Reaction is the same, so remove it
			err = app.reactions.Delete(existingReaction.ID)
			if err != nil {
				// Use your custom error message for deletion errors
				http.Error(w, fmt.Sprintf(ErrorMsgs().Delete, "storeReaction", err), http.StatusInternalServerError)
				log.Printf("Error deleting reaction: %v", err)
				return
			}
			// Send response back after successful deletion
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction removed"})
			if err != nil {
				return
			}
			return
		}

		// Otherwise, update the existing reaction
		err = app.reactions.Update(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
		if err != nil {
			// Use your custom error message for update errors
			http.Error(w, fmt.Sprintf(ErrorMsgs().Update, "storeReaction", err), http.StatusInternalServerError)
			log.Printf("Error updating reaction: %v", err)
			return
		}
		// Send response back after successful update
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction updated"})
		if err != nil {
			return
		}
		return
	}

	// If no existing reaction, insert a new reaction
	err = app.reactions.Upsert(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
	if err != nil {
		// Use your custom error message for insertion errors
		http.Error(w, fmt.Sprintf(ErrorMsgs().Insert, "storeReaction", err), http.StatusInternalServerError)
		log.Printf("Error inserting reaction: %v", err)
		return
	}

	// Check reaction and store it in the database, or handle errors
	// Respond with a JSON response
	w.Header().Set("Content-Type", "application/json")

	// Send a response indicating success
	//w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction added (in go)"})

	if err != nil {
		log.Printf(ErrorMsgs().Post, err)
		http.Error(w, err.Error(), 500)
		return
	}

	//http.Redirect(w, r, "/", http.StatusFound)
}
