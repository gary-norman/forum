package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	"strings"
	"time"

	"github.com/gary-norman/forum/internal/models"
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

// TODO use interface to get and return any type
func getRandomChannel(channelSlice []models.ChannelWithDaysAgo) models.ChannelWithDaysAgo {
	rndInt := rand.IntN(len(channelSlice))
	channel := channelSlice[rndInt]
	return channel
}

func getRandomUser(userSlice []models.User) models.User {
	rndInt := rand.IntN(len(userSlice))
	user := userSlice[rndInt]
	return user
}

func (app *app) register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("register_user")
	email := r.FormValue("register_email")
	validEmail, _ := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, email)
	password := r.FormValue("register_password")
	if len(username) < 5 || len(username) > 16 {
		w.WriteHeader(http.StatusNotAcceptable)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    http.StatusNotAcceptable,
			"message": "username must be between 5 and 16 characters",
		})
		if err != nil {
			log.Printf(ErrorMsgs().Encode, "register: username", err)
			return
		}
		return
	}
	if isValidPassword(password) != true {
		w.WriteHeader(http.StatusNotAcceptable)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
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
	if validEmail != true {
		w.WriteHeader(http.StatusNotAcceptable)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
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
	if emailExists == true {
		w.WriteHeader(http.StatusConflict)
		encErr := json.NewEncoder(w).Encode(map[string]interface{}{
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
	if userExists == true {
		w.WriteHeader(http.StatusConflict)
		encErr := json.NewEncoder(w).Encode(map[string]interface{}{
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
		encErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": errors.New(fmt.Sprintf(ErrorMsgs().Register, insertUser)),
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
	encErr := json.NewEncoder(w).Encode(map[string]interface{}{
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
		encErr := json.NewEncoder(w).Encode(map[string]interface{}{
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
			encErr := json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "Failed to create cookies",
				"body":    errors.New(fmt.Sprintf(ErrorMsgs().Cookies, "create", createCookiErr)),
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
		encErr := json.NewEncoder(w).Encode(map[string]interface{}{
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
		encErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "incorrect password",
		})
		if encErr != nil {
			log.Printf(ErrorMsgs().Encode, "login: fail", encErr)
			return
		}
	}
}

func (app *app) JoinedByCurrentUser(memberships []models.Membership) ([]models.Channel, error) {
	fmt.Println(Colors().Orange + "Checking if this user is a member of this channel" + Colors().Reset)
	fmt.Printf(ErrorMsgs().Divider)
	var channels []models.Channel
	for _, membership := range memberships {
		channel, err := app.channels.OwnedOrJoinedByCurrentUser(membership.ChannelID, "ID")
		if err != nil {
			return nil, errors.New(fmt.Sprintf(ErrorMsgs().KeyValuePair, "Error calling JoinedByCurrentUser > OwnedOrJoinedByCurrentUser", err))
		}
		channels = append(channels, channel[0])
	}
	if len(channels) > 0 {
		fmt.Println(Colors().Green + "Current user is a member of this channel" + Colors().Reset)
	} else {
		fmt.Println(Colors().Red + "Current user is not a member of this channel" + Colors().Reset)
	}
	return channels, nil
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
	log.Printf(Colors.Green+"%v logged out successfully!", user)
	encErr := json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Logged out successfully!",
	})
	if encErr != nil {
		log.Printf(ErrorMsgs().Encode, "logout: success", encErr)
		return
	}
}

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

func (app *app) getHome(w http.ResponseWriter, r *http.Request) {
	//var userLoggedIn bool
	userLoggedIn := true
	// SECTION --- posts and comments ---
	posts, postsErr := app.posts.All()
	if postsErr != nil {
		http.Error(w, postsErr.Error(), 500)
		return
	}
	// Retrieve total likes and dislikes for each post
	for i, post := range posts {
		likes, dislikes, likesErr := app.reactions.CountReactions(post.ID, 0) // Pass 0 for CommentID if it's a post
		//fmt.Printf("PostID: %v, Likes: %v, Disikes: %v\n", posts[i].ID, likes, dislikes)
		if likesErr != nil {
			log.Printf("Error counting reactions: %v", likesErr)
			likes, dislikes = 0, 0 // Default values if there is an error
		}
		posts[i].Likes += likes
		posts[i].Dislikes += dislikes
	}

	comments, commentsErr := app.comments.All()
	if commentsErr != nil {
		//http.Error(w, commentsErr.Error(), 500)
		log.Printf("Error counting comments: %v", commentsErr)
	}

	// Retrieve total likes and dislikes for each post
	for i, comment := range comments {
		likes, dislikes, likesErr := app.reactions.CountReactions(0, comment.ID) // Pass 0 for CommentID if it's a post
		//fmt.Printf("PostID: %v, Likes: %v, Disikes: %v\n", posts[i].ID, likes, dislikes)
		if likesErr != nil {
			log.Printf("Error counting reactions: %v", likesErr)
			likes, dislikes = 0, 0 // Default values if there is an error
		}
		comments[i].Likes += likes
		comments[i].Dislikes += dislikes
	}

	commentsWithDaysAgo := make([]models.CommentWithDaysAgo, len(comments))
	for index, comment := range comments {
		commentsWithDaysAgo[index] = models.CommentWithDaysAgo{
			Comment:   comment,
			TimeSince: getTimeSince(comment.Created),
			Replies:   []models.CommentWithDaysAgo{}, // Initialize with no replies
		}
	}
	// TODO make a function with interface to unify this and other withdaysago
	postsWithDaysAgo := make([]models.PostWithDaysAgo, len(posts))
	for index, post := range posts {
		/// Filter comments that belong to the current post based on the postID and CommentedPostID
		var postComments []models.CommentWithDaysAgo
		for _, comment := range commentsWithDaysAgo {
			// Match the postID with CommentedPostID
			if comment.Comment.CommentedPostID != nil && *comment.Comment.CommentedPostID == post.ID {
				// For each comment, recursively assign its replies
				commentWithReplies := getRepliesForComment(comment, commentsWithDaysAgo)
				postComments = append(postComments, commentWithReplies)
			}

		}
		postsWithDaysAgo[index] = models.PostWithDaysAgo{
			Post:      post,
			TimeSince: getTimeSince(post.Created),
			Comments:  postComments, // Only include comments related to the current post
		}
	}

	// SECTION --- user ---
	allUsers, allUsersErr := app.users.All()
	if allUsersErr != nil {
		log.Printf(ErrorMsgs().Query, "getHome> users > All", allUsersErr)
	}

	//attach following/follower numbers to each user
	for _, user := range allUsers {
		user.Followers, user.Following, allUsersErr = app.loyalty.CountUsers(user.ID)
	}

	randomUser := getRandomUser(allUsers)
	currentUser, currentUserErr := app.GetLoggedInUser(r)
	if currentUserErr != nil {
		log.Printf(ErrorMsgs().NotFound, "user", "current user", "GetLoggedInUser", currentUserErr)
		log.Printf(ErrorMsgs().KeyValuePair, "Current user", currentUser)
		userLoggedIn = false
	}

	//attach following/follower numbers to a random user
	randomUser.Followers, randomUser.Following, currentUserErr = app.loyalty.CountUsers(randomUser.ID)

	//attach following/follower numbers to currently logged-in user
	currentUser.Followers, currentUser.Following, currentUserErr = app.loyalty.CountUsers(currentUser.ID)

	//validTokens := app.cookies.QueryCookies(w, r, currentUser)
	//if validTokens == true {
	//	userLoggedIn = true
	//}
	currentUserName := "nouser"
	var currentUserID int
	var currentUserAvatar string
	var currentUserBio string

	// SECTION --- channels --
	allChannels, allChanErr := app.channels.All()
	if allChanErr != nil {
		log.Printf(ErrorMsgs().Query, "channels.All", allChanErr)
	}
	channelsWithDaysAgo := make([]models.ChannelWithDaysAgo, len(allChannels))
	for index, channel := range allChannels {
		channelsWithDaysAgo[index] = models.ChannelWithDaysAgo{
			Channel:   channel,
			TimeSince: getTimeSince(channel.Created),
		}
	}
	var thisChannel models.ChannelWithDaysAgo
	channelId, err := strconv.Atoi(r.PathValue("channelId"))
	foundChannels, err := app.channels.Search("id", channelId)
	if err != nil {
		fmt.Printf(ErrorMsgs().KeyValuePair, "getHome > found channels", err)
	}
	if len(foundChannels) > 0 {
		thisChannel = models.ChannelWithDaysAgo{
			Channel:   foundChannels[0],
			TimeSince: getTimeSince(foundChannels[0].Created),
		}
	} else {
		fmt.Printf(ErrorMsgs().KeyValuePair, "no channel found", "getting random channel")
		thisChannel = getRandomChannel(channelsWithDaysAgo)
	}

	thisChannelOwnerName, ownerErr := app.users.GetSingleUserValue(thisChannel.OwnerID, "ID", "username")
	if ownerErr != nil {
		log.Printf(ErrorMsgs().Query, "getHome > GetSingleUserValue", ownerErr)
	}
	var ownedChannels []models.Channel
	var joinedChannels []models.Channel
	var ownedAndJoinedChannels []models.Channel
	isJoinedOrOwned := false
	isOwned := false

	if userLoggedIn == true {
		currentUser.TimeSince = getTimeSince(currentUser.Created)
		currentUserName = currentUser.Username
		currentUserID = currentUser.ID
		currentUserAvatar = currentUser.Avatar
		currentUserBio = currentUser.Description
		// get owned and joined channels of current user
		memberships, memberErr := app.memberships.UserMemberships(currentUser.ID)
		if memberErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > UserMemberships", memberErr)
		}
		ownedChannels, ownedChannelsErr := app.channels.OwnedOrJoinedByCurrentUser(currentUser.ID, "OwnerID")
		if ownedChannelsErr != nil {
			log.Printf(ErrorMsgs().Query, "user owned channels", ownedChannelsErr)
		}
		joinedChannels, joinedChannelsErr := app.JoinedByCurrentUser(memberships)
		if joinedChannelsErr != nil {
			log.Printf(ErrorMsgs().Query, "user joined channels", joinedChannelsErr)
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
		log.Printf(ErrorMsgs().KeyValuePair, "getHome > AllForChannel", err)
	}
	// TODO get channel moderators

	// TODO get channel posts
	// TODO create structs for each type (ChannelTemplate etc) and serve them in the template struct
	// SECTION -- template ---
	templateData := models.TemplateData{
		AllUsers:    allUsers,
		RandomUser:  randomUser,
		CurrentUser: currentUser,
		//TODO get these values dynamically (NIL pointer reference)
		CurrentUserID:              currentUserID,
		CurrentUserName:            currentUserName,
		AllChannels:                allChannels,
		ThisChannel:                thisChannel,
		ThisChannelOwnerName:       thisChannelOwnerName,
		OwnedChannels:              ownedChannels,
		JoinedChannels:             joinedChannels,
		OwnedAndJoinedChannels:     ownedAndJoinedChannels,
		ThisChannelIsOwned:         isOwned,
		ThisChannelIsOwnedOrJoined: isJoinedOrOwned,
		ThisChannelRules:           thisChannelRules,
		Posts:                      postsWithDaysAgo,
		Avatar:                     currentUserAvatar,
		Bio:                        currentUserBio,
		Images:                     nil,
		Reactions:                  nil,
		NotifyPlaceholder: models.NotifyPlaceholder{
			Register: "",
			Login:    "",
		},
	}
	//models.JsonError(templateData)
	tpl, err := GetTemplate()
	if err != nil {
		log.Printf(ErrorMsgs().Parse, "templates", "getHome", err)
		return
	}

	t := tpl.Lookup("index.html")
	//if t == nil {
	//	log.Printf("Template not found: index.html")
	//	return
	//}

	if t == nil {
		http.Error(w, "Template is not initialized", http.StatusInternalServerError)
		return
	}

	data := templateData
	execErr := t.Execute(w, data)
	if execErr != nil {
		log.Printf(ErrorMsgs().Execute, execErr)
		return
	}
}

// Recursively fetch replies for each comment
func getRepliesForComment(comment models.CommentWithDaysAgo, commentsWithDaysAgo []models.CommentWithDaysAgo) models.CommentWithDaysAgo {
	// Find replies to the current comment
	var replies []models.CommentWithDaysAgo
	for _, possibleReply := range commentsWithDaysAgo {
		if possibleReply.Comment.CommentedCommentID != nil && *possibleReply.Comment.CommentedCommentID == comment.Comment.ID {
			replyWithReplies := getRepliesForComment(possibleReply, commentsWithDaysAgo) // Recursively get replies for this reply
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

func (app *app) storeChannel(w http.ResponseWriter, r *http.Request) {
	user, getUserErr := app.GetLoggedInUser(r)
	if getUserErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
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
		IsFlagged:   false,
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
		createChannelData.IsFlagged,
		createChannelData.IsMuted,
	)

	if insertErr != nil {
		log.Printf(ErrorMsgs().Post, insertErr)
		http.Error(w, insertErr.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *app) storeMembership(w http.ResponseWriter, r *http.Request) {
	user, getUserErr := app.GetLoggedInUser(r)
	if getUserErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if parseErr := r.ParseForm(); parseErr != nil {
		http.Error(w, parseErr.Error(), 400)
		log.Printf(ErrorMsgs().Parse, "storeMembership", parseErr)
		return
	}
	fmt.Printf("user: %v", user.Username)
	// get channelID
	joinedChannelID, convErr := strconv.Atoi(r.PostForm.Get("channelId"))
	if convErr != nil {
		log.Printf(ErrorMsgs().Convert, r.PostForm.Get("channelId"), "StoreMembership > GetChannelID", convErr)
		log.Printf("Unable to convert %v to integer\n", r.PostForm.Get("channelId"))
	}
	// get slice of channels (in this case it is only 1, but the function still returns a slice)
	channels, err := app.channels.Search("id", joinedChannelID)
	if err != nil {
		log.Printf(ErrorMsgs().Query, "channel", err)
	}
	// get the channel object
	channel := channels[0]

	createMembershipData := models.Membership{
		UserID:    user.ID,
		ChannelID: joinedChannelID,
	}
	// send memberships struct to DB
	insertErr := app.memberships.Insert(
		createMembershipData.UserID,
		createMembershipData.ChannelID,
	)
	if insertErr != nil {
		log.Printf(ErrorMsgs().Post, insertErr)
		http.Error(w, insertErr.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encErr := json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"message": fmt.Sprintf("Welcome to %v!", channel.Name),
	})
	if encErr != nil {
		log.Printf(ErrorMsgs().Encode, "storeMembership: Accepted", encErr)
		return
	}
}

func (app *app) requestMembership(w http.ResponseWriter, r *http.Request, userID, channelID int) {

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
		//TODO expect: 100-continue on the request header (for all of these)
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
		// Use your custom error message for insertion errors
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

func (app *app) CreateAndInsertRule(w http.ResponseWriter, r *http.Request) {
	channelId, err := strconv.Atoi(r.PathValue("channelId"))
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "CreateAndInsertRule > convert channelId to int", err)
	}

	// Get the "rules" input value
	rulesJSON := r.FormValue("rules")
	if rulesJSON == "" { // TODO send this message to the user
		log.Printf(ErrorMsgs().KeyValuePair, "message to user", "you have not added or removed any rules")
	}

	// Decode JSON into a slice of Rule structs
	var rules []models.POSTRule
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	for _, rule := range rules {
		id, found := strings.CutPrefix(rule.ID, "existing-channel-rule-")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "CreateAndInsertRule > id => idInt", err)
		}
		if found == true {
			err := app.rules.DeleteRule(channelId, idInt)
			if err != nil {
				log.Printf(ErrorMsgs().KeyValuePair, "CreateAndInsertRule > DeleteRule", err)
			}
		} else {
			ruleId, err := app.rules.CreateRule(rule.Text)
			if err != nil {
				log.Printf(ErrorMsgs().KeyValuePair, "CreateAndInsertRule > CreateRule", err)
			}
			err = app.rules.InsertRule(channelId, ruleId)
			if err != nil {
				log.Printf(ErrorMsgs().KeyValuePair, "CreateAndInsertRule > InsertRule", err)
			}
		}
	}
	// TODO redirect to /channels/{{.channelID}}
	http.Redirect(w, r, "/channels/"+r.PathValue("channelId"), http.StatusFound)
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
