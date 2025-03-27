package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (app *app) getThisUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userLoggedIn := true
	currentUser, ok := getUserFromContext(r.Context())
	if !ok {
		fmt.Printf(ErrorMsgs().NotFound, "currentUser", "getThisUser", "_")
		userLoggedIn = false
	}

	userId, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		fmt.Printf(ErrorMsgs().KeyValuePair, "convert userID", err)
	}
	thisUser, err := app.users.GetUserByID(userId)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "getHome > thisUser", err)
	}
	thisUser.Followers, thisUser.Following, err = app.loyalty.CountUsers(thisUser.ID)
	if err != nil {
		fmt.Printf(ErrorMsgs().KeyValuePair, "getHome > thisUser loyalty", err)
	}

	if userLoggedIn {
		currentUser.Followers, currentUser.Following, err = app.loyalty.CountUsers(currentUser.ID)
		if err != nil {
			fmt.Printf(ErrorMsgs().KeyValuePair, "getHome > currentUser loyalty", err)
		}
	}

	TemplateData.ThisUser = thisUser
	TemplateData.CurrentUser = currentUser
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
