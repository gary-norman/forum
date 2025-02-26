package main

import (
	"encoding/json"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"html/template"
	"io"
	"log"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var TemplateData models.TemplateData

// SECTION ------- template handlers ----------

func (app *app) getHome(w http.ResponseWriter, r *http.Request) {
	// SECTION --- posts and comments ---
	// var userLoggedIn bool
	userLoggedIn := true
	allPosts, err := app.posts.All()
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error fetching all posts", err)
	}

	allComments, err := app.comments.All()
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error fetching all comments", err)
	}
	// Retrieve total likes and dislikes for each post
	allPosts = app.getPostsLikesAndDislikes(allPosts)
	// Retrieve total likes and dislikes for each comment
	allComments = app.getCommentsLikesAndDislikes(allComments)

	for p := range allPosts {
		models.UpdateTimeSince(&allPosts[p])
	}
	for c := range allComments {
		models.UpdateTimeSince(&allComments[c])
	}

	allPosts = app.getPostsComments(allPosts, allComments)

	// SECTION --- user ---
	allUsers, allUsersErr := app.users.All()
	if allUsersErr != nil {
		log.Printf(ErrorMsgs().Query, "getHome> users > All", allUsersErr)
	}
	for u := range allUsers {
		models.UpdateTimeSince(&allUsers[u])
	}

	// attach following/follower numbers to each user
	for _, user := range allUsers {
		user.Followers, user.Following, allUsersErr = app.loyalty.CountUsers(user.ID)
		if allUsersErr != nil {
			log.Printf(ErrorMsgs().Query, "getHome> users > All > loyalty", allUsersErr)
		}
	}

	randomUser := getRandomUser(allUsers)
	currentUser, currentUserErr := app.GetLoggedInUser(r)
	if currentUserErr != nil {
		log.Printf(ErrorMsgs().NotFound, "user", "current user", "GetLoggedInUser", currentUserErr)
		log.Printf(ErrorMsgs().KeyValuePair, "Current user", currentUser)
		userLoggedIn = false
	}

	// attach following/follower numbers to the random user
	randomUser.Followers, randomUser.Following, currentUserErr = app.loyalty.CountUsers(randomUser.ID)
	if currentUserErr != nil {
		log.Printf(ErrorMsgs().Query, "getHome> users > All", allUsersErr)
	}

	//validTokens := app.cookies.QueryCookies(w, r, currentUser)
	//if validTokens == true {
	//	userLoggedIn = true
	//}

	// SECTION --- channels --
	allChannels, err := app.channels.All()
	if err != nil {
		log.Printf(ErrorMsgs().Query, "channels.All", err)
	}
	for c := range allChannels {
		models.UpdateTimeSince(&allChannels[c])
	}
	var thisChannel models.Channel
	channelId, err := strconv.Atoi(r.PathValue("channelId"))
	if err != nil {
		fmt.Printf(ErrorMsgs().KeyValuePair, "convert channelId", err)
	}
	foundChannels, err := app.channels.Search("id", channelId)
	if err != nil {
		fmt.Printf(ErrorMsgs().KeyValuePair, "getHome > found channels", err)
	}
	if len(foundChannels) > 0 {
		thisChannel = foundChannels[0]
	} else {
		fmt.Printf(ErrorMsgs().KeyValuePair, "no channel found", "getting random channel")
		thisChannel = getRandomChannel(allChannels)
	}
	thisChannelOwnerName, err := app.users.GetSingleUserValue(thisChannel.OwnerID, "ID", "username")
	if err != nil {
		log.Printf(ErrorMsgs().Query, "getHome > GetSingleUserValue", err)
	}
	var ownedChannels []models.Channel
	var joinedChannels []models.Channel
	var ownedAndJoinedChannels []models.Channel
	isJoinedOrOwned := false
	isOwned := false

	if userLoggedIn {
		// attach following/follower numbers to currently logged-in user
		currentUser.Followers, currentUser.Following, err = app.loyalty.CountUsers(currentUser.ID)
		if err != nil {
			fmt.Printf(ErrorMsgs().KeyValuePair, "getHome > currentUser loyalty", err)
		}
		// get owned and joined channels of current user
		memberships, memberErr := app.memberships.UserMemberships(currentUser.ID)
		if memberErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > UserMemberships", memberErr)
		}
		ownedChannels, err = app.channels.OwnedOrJoinedByCurrentUser(currentUser.ID, "OwnerID")
		if err != nil {
			log.Printf(ErrorMsgs().Query, "user owned channels", err)
		}
		joinedChannels, err = app.JoinedByCurrentUser(memberships)
		if err != nil {
			log.Printf(ErrorMsgs().Query, "user joined channels", err)
		}
		ownedAndJoinedChannels = append(ownedChannels, joinedChannels...)
		isOwned = currentUser.ID == thisChannel.OwnerID
		joined := false
		for _, channel := range joinedChannels {
			if thisChannel.ID == channel.ID {
				joined = true
				break
			}
		}
		isJoinedOrOwned = isOwned || joined
	}
	// get all rules for the current channel
	thisChannelRules, err := app.rules.AllForChannel(thisChannel.ID)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "getHome > rules.AllForChannel", err)
	}
	// TODO get channel moderators

	thisChannelPosts, err := app.posts.GetPostsByChannel(thisChannel.ID)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "getHome > allPosts.GetPostsByChannel", err)
	}
	fmt.Printf(ErrorMsgs().KeyValuePair, "thisChannelPosts", len(thisChannelPosts))
	// Retrieve total likes and dislikes for each Channel post
	thisChannelPosts = app.getPostsLikesAndDislikes(thisChannelPosts)
	app.getPostsComments(thisChannelPosts, allComments)

	// SECTION -- template ---

	// ---------- users ----------
	TemplateData.AllUsers = allUsers
	TemplateData.RandomUser = randomUser
	TemplateData.CurrentUser = currentUser
	// ---------- allPosts ----------
	TemplateData.Posts = allPosts
	// ---------- channels ----------
	TemplateData.AllChannels = allChannels
	TemplateData.ThisChannel = thisChannel
	TemplateData.ThisChannelOwnerName = thisChannelOwnerName
	TemplateData.OwnedChannels = ownedChannels
	TemplateData.JoinedChannels = joinedChannels
	TemplateData.OwnedAndJoinedChannels = ownedAndJoinedChannels
	TemplateData.ThisChannelIsOwned = isOwned
	TemplateData.ThisChannelIsOwnedOrJoined = isJoinedOrOwned
	TemplateData.ThisChannelRules = thisChannelRules
	TemplateData.ThisChannelPosts = thisChannelPosts
	// ---------- misc ----------
	TemplateData.Images = nil
	TemplateData.Reactions = nil

	// models.JsonError(TemplateData)
	tpl, err := GetTemplate()
	if err != nil {
		log.Printf(ErrorMsgs().Parse, "templates", "getHome", err)
		return
	}

	t := tpl.Lookup("index.html")

	if t == nil {
		http.Error(w, "Template is not initialized", http.StatusInternalServerError)
		return
	}

	execErr := t.Execute(w, TemplateData)
	if execErr != nil {
		log.Printf(ErrorMsgs().Execute, execErr)
		return
	}
}

// SECTION ------- user login handlers ----------

func (app *app) register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("register_user")
	email := r.FormValue("register_email")
	validEmail, _ := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, email)
	password := r.FormValue("register_password")
	if len(username) < 5 || len(username) > 16 {
		w.WriteHeader(http.StatusNotAcceptable)
		err := json.NewEncoder(w).Encode(map[string]any{
			"code":    http.StatusNotAcceptable,
			"message": "username must be between 5 and 16 characters",
		})
		if err != nil {
			log.Printf(ErrorMsgs().Encode, "register: username", err)
			return
		}
		return
	}
	if !isValidPassword(password) {
		w.WriteHeader(http.StatusNotAcceptable)
		err := json.NewEncoder(w).Encode(map[string]any{
			"code": http.StatusNotAcceptable,
			"message": "password must contain at least one number and one uppercase and lowercase letter," +
				"and at least 8 or more characters",
		})
		if err != nil {
			log.Printf(ErrorMsgs().Encode, "register: password", err)
			return
		}
		return
	}
	if !validEmail {
		w.WriteHeader(http.StatusNotAcceptable)
		err := json.NewEncoder(w).Encode(map[string]any{
			"code":    http.StatusNotAcceptable,
			"message": "please enter a valid email address",
		})
		if err != nil {
			log.Printf(ErrorMsgs().Encode, "register: validEmail", err)
			return
		}
		return
	}
	emailExists, err := app.users.QueryUserEmailExists(email)
	if emailExists {
		w.WriteHeader(http.StatusConflict)
		encErr := json.NewEncoder(w).Encode(map[string]any{
			"code":    http.StatusConflict,
			"message": "an account is already registered to that email address",
			"body":    err,
		})
		if encErr != nil {
			log.Printf(ErrorMsgs().Encode, "register: emailExists", encErr)
			return
		}
		return
	}
	userExists, err := app.users.QueryUserNameExists(username)
	if userExists {
		w.WriteHeader(http.StatusConflict)
		encErr := json.NewEncoder(w).Encode(map[string]any{
			"code":    http.StatusConflict,
			"message": "an account is already registered to that username",
			"body":    err,
		})
		if encErr != nil {
			log.Printf(ErrorMsgs().Encode, "register: userExists", encErr)
			return
		}
		return
	}
	hashedPassword, _ := models.HashPassword(password)
	insertUser := app.users.Insert(
		username,
		email,
		hashedPassword,
		"",
		"",
		"noimage",
		"default.png",
		"")

	if insertUser != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encErr := json.NewEncoder(w).Encode(map[string]any{
			"code":    http.StatusInternalServerError,
			"message": fmt.Errorf(ErrorMsgs().Register, insertUser),
		})
		if encErr != nil {
			log.Printf(ErrorMsgs().Encode, "register: insertErr", encErr)
			return
		}
		return
	}
	type FormFields struct {
		Fields map[string][]string `json:"formValues"`
	}
	formFields := make(map[string][]string)
	for field, value := range r.Form {
		fieldName := field
		formFields[fieldName] = append(formFields[fieldName], value...)
	}
	// Send success response
	w.WriteHeader(http.StatusOK)
	encErr := json.NewEncoder(w).Encode(map[string]any{
		"code":    http.StatusOK,
		"message": "Registration successful!",
		"body":    FormFields{Fields: formFields},
	})
	if encErr != nil {
		log.Printf(ErrorMsgs().Encode, "register: send success", encErr)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	Colors := models.CreateColors()
	// Parse JSON from the request body
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Printf("Failed to parse request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	login := credentials.Username
	password := credentials.Password
	fmt.Printf(Colors.Orange+"Attempting login for "+Colors.White+"%v\n"+Colors.Reset, login)
	fmt.Printf(ErrorMsgs().Divider)

	user, getUserErr := app.users.GetUserFromLogin(login, "login")
	if getUserErr != nil {
		log.Printf(ErrorMsgs().NotFound, "either", login, "login > GetUserFromLogin", getUserErr)
		w.WriteHeader(http.StatusUnauthorized)
		encErr := json.NewEncoder(w).Encode(map[string]any{
			"code":    http.StatusUnauthorized,
			"message": "User not found",
		})
		if encErr != nil {
			log.Printf(ErrorMsgs().Encode, "login: CreateCookies", encErr)
			return
		}
		return
	}

	if models.CheckPasswordHash(password, user.HashedPassword) {
		fmt.Printf(Colors.Green+"Passwords for %v match\n"+Colors.Reset, user.Username)
		// Set Session Token and CSRF Token cookies
		createCookiErr := app.cookies.CreateCookies(w, user)
		if createCookiErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encErr := json.NewEncoder(w).Encode(map[string]any{
				"code":    http.StatusInternalServerError,
				"message": "Failed to create cookies",
				"body":    fmt.Errorf(ErrorMsgs().Cookies, "create", createCookiErr),
			})
			if encErr != nil {
				log.Printf(ErrorMsgs().Encode, "login: CreateCookies", encErr)
				return
			}
			return
		}
		// Respond with a successful login message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encErr := json.NewEncoder(w).Encode(map[string]any{
			"code":    http.StatusOK,
			"message": fmt.Sprintf("Welcome, %s! Login successful.", user.Username),
		})
		if encErr != nil {
			log.Printf(ErrorMsgs().Encode, "login: success", encErr)
			return
		}
	} else {
		// Respond with an unsuccessful login message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encErr := json.NewEncoder(w).Encode(map[string]any{
			"code":    http.StatusUnauthorized,
			"message": "incorrect password",
		})
		if encErr != nil {
			log.Printf(ErrorMsgs().Encode, "login: fail", encErr)
			return
		}
	}
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

	// Delete the Session Token and CSRF Token cookies
	delCookiErr := app.cookies.DeleteCookies(w, user)
	if delCookiErr != nil {
		log.Printf(ErrorMsgs().Cookies, "delete", delCookiErr)
	}
	// send user confirmation
	log.Printf(Colors.Green+"%v logged out successfully!", user.Username)
	encErr := json.NewEncoder(w).Encode(map[string]any{
		"code":    http.StatusOK,
		"message": "Logged out successfully!",
	})
	if encErr != nil {
		log.Printf(ErrorMsgs().Encode, "logout: success", encErr)
		return
	}
}

// SECTION ------- routing handlers ----------
func (app *app) protected(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("username")
	var user *models.User
	user, getUserErr := app.users.GetUserFromLogin(login, "protected")
	if getUserErr != nil {
		log.Printf("protected route for %v failed with error: %v", login, getUserErr)
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

// SECTION ------- reactions handlers ----------

//type reactions interface {
//	Likes(count int)
//	Dislikes(count int)
//}

// getPostsLikesAndDislikes updates the reactions of each post in the given slice
func (app *app) getPostsLikesAndDislikes(posts []models.Post) []models.Post {
	for i, post := range posts {
		likes, dislikes, err := app.reactions.CountReactions(post.ID, 0) // Pass 0 for CommentID if it's a post
		// fmt.Printf("PostID: %v, Likes: %v, Dislikes: %v\n", posts[i].ID, likes, dislikes)
		if err != nil {
			log.Printf("Error counting reactions: %v", err)
			likes, dislikes = 0, 0 // Default values if there is an error
		}
		posts[i].React(likes, dislikes)
	}
	return posts
}

// getCommentsLikesAndDislikes updates the reactions of each comment in the given slice
func (app *app) getCommentsLikesAndDislikes(comments []models.Comment) []models.Comment {
	for i, comment := range comments {
		likes, dislikes, likesErr := app.reactions.CountReactions(0, comment.ID) // Pass 0 for PostID if it's a comment
		// fmt.Printf("PostID: %v, Likes: %v, Dislikes: %v\n", posts[i].ID, likes, dislikes)
		if likesErr != nil {
			log.Printf("Error counting reactions: %v", likesErr)
			likes, dislikes = 0, 0 // Default values if there is an error
		}
		comments[i].React(likes, dislikes)
	}
	return comments
}

// SECTION ------- user handlers ----------
func (app *app) editUserDetails(w http.ResponseWriter, r *http.Request) {
	user, err := app.GetLoggedInUser(r)
	if err != nil {
		log.Printf(ErrorMsgs().NotFound, "user", "current user", "GetLoggedInUser", err)
		return
	}
	if user.Username == "" {
		log.Printf(ErrorMsgs().NotFound, "user", "current user", "GetLoggedInUser", err)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Printf(ErrorMsgs().Parse, "editUserDetails", err)
		return
	}
	currentAvatar := user.Avatar
	fmt.Printf("currentAvatar: %v", currentAvatar)
	user.Avatar = GetFileName(r, "file-drop", "editUserDetails", "user")
	if currentAvatar != "noimage" {
		user.Avatar = currentAvatar
	}
	currentDescription := r.FormValue("bio")
	if currentDescription != "" {
		user.Description = currentDescription
	}
	currentName := r.FormValue("name")
	if currentName != "" {
		user.Username = currentName
	}
	editErr := app.users.Edit(user)
	if editErr != nil {
		log.Printf(ErrorMsgs().Edit, user.Username, "EditUserDetails", editErr)
	}
	if err := app.cookies.CreateCookies(w, user); err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "error creating cookies", err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// SECTION ------- post handlers ----------
func (app *app) createPost(w http.ResponseWriter, r *http.Request) {
	tpl, parseErr := template.ParseFiles("./assets/templates/posts.create.html")
	if parseErr != nil {
		http.Error(w, parseErr.Error(), 500)
		log.Printf(ErrorMsgs().Parse, "./assets/templates/posts.create.html", "createPost", parseErr)
		return
	}

	execErr := tpl.Execute(w, nil)
	if execErr != nil {
		log.Printf(ErrorMsgs().Execute, execErr)
		return
	}
}

func (app *app) storePost(w http.ResponseWriter, r *http.Request) {
	user, getUserErr := app.GetLoggedInUser(r)
	if getUserErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
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

func (app *app) storeReaction(w http.ResponseWriter, r *http.Request) {
	// log.Printf("using storeReaction()")

	// Variable to hold the decoded data
	var reactionData models.Reaction

	// Decode the JSON request body into the reactionData variable
	err := json.NewDecoder(r.Body).Decode(&reactionData)
	fmt.Printf("reactionData received: %v\n", &reactionData)
	if err != nil {
		// Use the error message from the Errors struct for decoding errors
		log.Println("Error decoding JSON")
		log.Printf(ErrorMsgs().Parse, "storeReaction", err)
		// TODO expect: 100-continue on the request header (for all of these)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//// Validate that at least one of reactedPostID or reactedCommentID is non-zero
	if (reactionData.ReactedPostID == nil || *reactionData.ReactedPostID == 0) && (reactionData.ReactedCommentID == nil || *reactionData.ReactedCommentID == 0) {
		log.Println("You must react to either a post or a comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postID, commentID := 0, 0

	if reactionData.ReactedPostID != nil {
		// log.Println("ReactedPostID:", *reactionData.ReactedPostID)
		postID = *reactionData.ReactedPostID
	}

	if reactionData.ReactedCommentID != nil {
		// log.Printf("ReactedCommentID: %d", *reactionData.ReactedPostID)
		commentID = *reactionData.ReactedCommentID
	}

	fmt.Printf("commentID after conversion: %v\n", commentID)
	fmt.Printf("postID after conversion: %v\n", postID)

	// Check if the user already reacted (like/dislike) and update or delete the reaction if needed
	existingReaction, err := app.reactions.CheckExistingReaction(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
	if err != nil {
		// Use your custom error message for fetching errors
		log.Printf(ErrorMsgs().Read, "storeReaction > app.reactions.CheckExistingReaction", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// If there is an existing reaction, toggle it (i.e., remove it if the user reacts again to the same thing)
	if existingReaction != nil {
		log.Printf("Existing Reaction: %+v", existingReaction)

		// If the user likes a post or comment again, remove the like/dislike (toggle behavior)
		if existingReaction.Liked == reactionData.Liked && existingReaction.Disliked == reactionData.Disliked {
			// Reaction is the same, so remove it
			err = app.reactions.Delete(existingReaction.ID)
			if err != nil {
				// Use your custom error message for deletion errors
				log.Printf(ErrorMsgs().Delete, "storeReaction > app.reactions.Delete", err)
				w.WriteHeader(http.StatusInternalServerError)
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
			log.Printf(ErrorMsgs().Delete, "storeReaction > app.reactions.Update", err)
			w.WriteHeader(http.StatusInternalServerError)
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
	fmt.Printf("Upserting liked: %v, disliked: %v, authorID: %v, postID: %v, commentID: %v\n", reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
	if err != nil {
		log.Printf(ErrorMsgs().Delete, "storeReaction > app.reactions.Upsert", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Check reaction and store it in the database, or handle errors
	// Respond with a JSON response
	w.Header().Set("Content-Type", "application/json")
	// Send a response indicating success
	// w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction added (in go)"})
	if err != nil {
		log.Printf(ErrorMsgs().Post, err)
		http.Error(w, err.Error(), 500)
		return
	}
	// http.Redirect(w, r, "/", http.StatusFound)
}

func (app *app) storeComment(w http.ResponseWriter, r *http.Request) {
	// SECTION getting user
	user, getUserErr := app.GetLoggedInUser(r)
	if getUserErr != nil {
		http.Error(w, getUserErr.Error(), http.StatusUnauthorized)
		return
	}

	// Check if the method is POST, otherwise return Method Not Allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// SECTION retrieving comment form data
	// Variable to hold the decoded data
	var commentData models.Comment

	parseErr := r.ParseMultipartForm(10 << 20)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), 400)
		log.Printf(ErrorMsgs().Parse, "storeComment", parseErr)
		return
	}

	// SECTION setting channel data
	// Get channel data
	selectionJSON := r.PostForm.Get("channel")
	if selectionJSON == "" {
		http.Error(w, "No selection provided for channel", http.StatusBadRequest)
		return
	}

	var channelData models.ChannelData
	if err := json.Unmarshal([]byte(selectionJSON), &channelData); err != nil {
		log.Printf(ErrorMsgs().Unmarshal, selectionJSON, err)
		http.Error(w, "Invalid selection format", http.StatusBadRequest)
		return
	}

	postIDStr := r.PostForm.Get("postID")
	commentIDStr := r.PostForm.Get("commentID")

	// Convert postIDStr to an integer
	postID, postConvErr := strconv.Atoi(postIDStr)
	if postConvErr != nil {
		log.Printf("Error converting postID: %v", postID)
	}

	// Convert commentIDStr to an integer
	commentID, commentConvErr := strconv.Atoi(commentIDStr)
	if commentConvErr != nil {
		log.Printf("Error converting commentID: %v", commentID)
	}

	// Assign the returned values
	commentData = models.Comment{
		Content:            r.PostForm.Get("content"),
		CommentedPostID:    &postID,
		CommentedCommentID: &commentID,
		IsCommentable:      true,
		IsFlagged:          false,
		Author:             user.Username,
		AuthorID:           user.ID,
		AuthorAvatar:       user.Avatar,
		ChannelName:        channelData.ChannelName,
		ChannelID:          0,
	}

	commentData.ChannelID, _ = strconv.Atoi(channelData.ChannelID)

	// Log the values
	fmt.Printf("commentData.CommentedPostID: %v\n", *commentData.CommentedPostID)
	fmt.Printf("commentData.CommentedCommentID: %v\n", *commentData.CommentedCommentID)

	// Insert the comment
	insertErr := app.comments.Upsert(commentData)

	if insertErr != nil {
		log.Printf(ErrorMsgs().Comment, insertErr)
		http.Error(w, insertErr.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *app) getPostsComments(posts []models.Post, comments []models.Comment) []models.Post {
	for p, post := range posts {
		/// Filter comments that belong to the current post based on the postID and CommentedPostID
		var postComments []models.Comment
		for c, comment := range comments {
			comments[c].UpdateTimeSince()
			// Match the postID with CommentedPostID
			if comment.CommentedPostID != nil && *comment.CommentedPostID == post.ID {
				// For each comment, recursively assign its replies
				commentWithReplies := getRepliesForComment(comment, comments)
				postComments = append(postComments, commentWithReplies)
			}
		}
		posts[p].Comments = postComments
	}
	return posts
}

// SECTION ------- channel handlers ----------

// SECTION ------- helper functions ----------
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
	return hasUpper
}

func getTimeSince(created time.Time) string {
	now := time.Now()
	hours := now.Sub(created).Hours()
	var timeSince string
	if hours > 24 {
		timeSince = fmt.Sprintf("%.0f days ago", hours/24)
	} else if hours > 1 {
		timeSince = fmt.Sprintf("%.0f hours ago", hours)
	} else if minutes := now.Sub(created).Minutes(); minutes > 1 {
		timeSince = fmt.Sprintf("%.0f minutes ago", minutes)
	} else {
		timeSince = "just now"
	}
	return timeSince
}

func getRandomChannel(channelSlice []models.Channel) models.Channel {
	rndInt := rand.IntN(len(channelSlice))
	channel := channelSlice[rndInt]
	return channel
}

func getRandomUser(userSlice []models.User) models.User {
	rndInt := rand.IntN(len(userSlice))
	user := userSlice[rndInt]
	return user
}

// getRepliesForComment Recursively fetches replies for each comment
func getRepliesForComment(comment models.Comment, comments []models.Comment) models.Comment {
	// Find replies to the current comment
	var replies []models.Comment
	for _, possibleReply := range comments {
		if possibleReply.CommentedCommentID != nil && *possibleReply.CommentedCommentID == comment.ID {
			replyWithReplies := getRepliesForComment(possibleReply, comments) // Recursively get replies for this reply
			replies = append(replies, replyWithReplies)
		}
	}
	// If no replies are found, we can avoid unnecessary recursion
	if len(replies) > 0 {
		comment.Replies = replies
	}
	// Return the comment with its replies
	return comment
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
	// base = base[:len(base)-len(ext)]

	// Create the new file name with the UUID and extension
	newFilePath := filepath.Join(filepath.Dir(oldFilePath), newUUID+ext)

	return newFilePath
}
