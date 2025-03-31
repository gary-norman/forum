package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gary-norman/forum/internal/models"
)

func (app *app) getThisUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var thisUser models.User

	// userLoggedIn := true
	// currentUser, ok := getUserFromContext(r.Context())
	// if !ok {
	// 	fmt.Printf(ErrorMsgs().NotFound, "currentUser", "getThisUser", "_")
	// 	userLoggedIn = false
	// }

	// Parse User ID from URL
	userId, err := models.GetIntFromPathValue(r.PathValue("userId"))
	if err != nil {
		http.Error(w, `{"error": "invalid user ID"}`, http.StatusBadRequest)
	}

	// Fetch the user
	user, err := app.users.GetUserByID(userId)
	if err != nil {
		http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
	}

	// Fetch user loyalty
	if err == nil {
		thisUser = user
		thisUser.Followers, thisUser.Following, err = app.loyalty.CountUsers(thisUser.ID)
		if err != nil {
			http.Error(w, `{"error": "error fetching user loyalty"}`, http.StatusInternalServerError)
		}
	}

	// if userLoggedIn {
	// 	currentUser.Followers, currentUser.Following, err = app.loyalty.CountUsers(currentUser.ID)
	// 	if err != nil {
	// 		fmt.Printf(ErrorMsgs().KeyValuePair, "getHome > currentUser loyalty", err)
	// 	}
	// }

	TemplateData.ThisUser = thisUser
	// TemplateData.CurrentUser = currentUser

	fmt.Printf(ErrorMsgs().KeyValuePair, "TemplateData.ThisUser", TemplateData.ThisUser.Username)

	response := map[string]any{
		"user":      thisUser.Username,
		"followers": thisUser.Followers,
		"following": thisUser.Following,
	}

	// Write the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "error encoding JSON"}`, http.StatusInternalServerError)
	}
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
