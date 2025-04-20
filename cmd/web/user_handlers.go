package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gary-norman/forum/internal/models"
)

func (app *app) getThisUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var thisUser models.User

	userLoggedIn := true
	currentUser, ok := getUserFromContext(r.Context())
	if !ok {
		log.Printf(ErrorMsgs().NotFound, "currentUser", "getThisUser", "_")
		userLoggedIn = false
	}

	// Parse User ID from URL
	userId, err := models.GetIntFromPathValue(r.PathValue("userId"))
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "error parsing user ID", err)
		// http.Error(w, `{"error": "invalid user ID"}`, http.StatusBadRequest)
	}

	// Fetch the user
	user, err := app.users.GetUserByID(userId)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "error fetching user", err)
		// http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
	}

	// Fetch user loyalty
	if err == nil {
		thisUser = user
		thisUser.Followers, thisUser.Following, err = app.loyalty.CountUsers(thisUser.ID)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "error fetching user loyalty", err)
			// http.Error(w, `{"error": "error fetching user loyalty"}`, http.StatusInternalServerError)
		}
	}

	// Fetch user posts
	posts, err := app.posts.GetPostsByUserID(thisUser.ID)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "error fetching user posts", err)
		// http.Error(w, `{"error": "error fetching user posts"}`, http.StatusInternalServerError)
	}

	// Fetch channel name for posts
	for p := range posts {
		posts[p].ChannelID, posts[p].ChannelName, err = app.GetChannelInfoFromPostID(posts[p].ID)
		if err != nil {
			http.Error(w, `{"error": "error fetching channel info"}`, http.StatusInternalServerError)
		}

		models.UpdateTimeSince(&posts[p])
	}

	if userLoggedIn {
		currentUser.Followers, currentUser.Following, err = app.loyalty.CountUsers(currentUser.ID)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > currentUser loyalty", err)
		}
	}

	models.UpdateTimeSince(&thisUser)

	data := models.UserPage{
		CurrentUser: currentUser,
		Instance:    "user-page",
		ThisUser:    thisUser,
		Posts:       posts,
		ImagePaths:  app.paths,
	}

	renderPageData(w, data)
}

func (app *app) editUserDetails(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		log.Printf(ErrorMsgs().KeyValuePair, "user not found in context", "editUserDetails")
		return
	}

	err := r.ParseMultipartForm(10 << 20)
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
