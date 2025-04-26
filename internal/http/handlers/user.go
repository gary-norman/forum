package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gary-norman/forum/internal/app"
	mw "github.com/gary-norman/forum/internal/http/middleware"
	"github.com/gary-norman/forum/internal/models"
	"github.com/gary-norman/forum/internal/view"
)

type UserHandler struct {
	App      *app.App
	Reaction *ReactionHandler
	Comment  *CommentHandler
	Channel  *ChannelHandler
}

func (u *UserHandler) GetThisUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var thisUser models.User

	userLoggedIn := true
	currentUser, ok := mw.GetUserFromContext(r.Context())
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
	user, err := u.App.Users.GetUserByID(userId)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "error fetching user", err)
		// http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
	}

	// Fetch user loyalty
	if err == nil {
		thisUser = user
		thisUser.Followers, thisUser.Following, err = u.App.Loyalty.CountUsers(thisUser.ID)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "error fetching user loyalty", err)
			// http.Error(w, `{"error": "error fetching user loyalty"}`, http.StatusInternalServerError)
		}
	}

	// Fetch user posts
	posts, err := u.App.Posts.GetPostsByUserID(thisUser.ID)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "error fetching user posts", err)
		// http.Error(w, `{"error": "error fetching user posts"}`, http.StatusInternalServerError)
	}

	// Fetch channel name for posts
	for p := range posts {
		posts[p].ChannelID, posts[p].ChannelName, err = u.Channel.GetChannelInfoFromPostID(posts[p].ID)
		if err != nil {
			http.Error(w, `{"error": "error fetching channel info"}`, http.StatusInternalServerError)
		}

		models.UpdateTimeSince(&posts[p])
	}

	if userLoggedIn {
		currentUser.Followers, currentUser.Following, err = u.App.Loyalty.CountUsers(currentUser.ID)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > currentUser loyalty", err)
		}
	}

	models.UpdateTimeSince(&thisUser)

	data := models.UserPage{
		UserID:      models.NewUUIDField(), // Default value of 0 for logged out users
		CurrentUser: currentUser,
		Instance:    "user-page",
		ThisUser:    thisUser,
		Posts:       posts,
		ImagePaths:  u.App.Paths,
	}

	view.RenderPageData(w, data)
}

// GetLoggedInUser gets the currently logged-in user from the session token and returns the user's struct
func (u *UserHandler) GetLoggedInUser(r *http.Request) (*models.User, error) {
	// Get the username from the request cookie
	userCookie, getCookieErr := r.Cookie("username")
	if getCookieErr != nil {
		log.Printf(ErrorMsgs().Cookies, "get", getCookieErr)
		return nil, getCookieErr
	}
	var username string
	if userCookie != nil {
		username = userCookie.Value
	}
	fmt.Printf(ErrorMsgs().KeyValuePair, "Username", username)
	if username == "" {
		return nil, errors.New("no user is logged in")
	}
	user, getUserErr := u.App.Users.GetUserByUsername(username, "GetLoggedInUser")
	if getUserErr != nil {
		return nil, getUserErr
	}
	return user, nil
}

func (u *UserHandler) EditUserDetails(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
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
	editErr := u.App.Users.Edit(user)
	if editErr != nil {
		log.Printf(ErrorMsgs().Edit, user.Username, "EditUserDetails", editErr)
	}
	if err := u.App.Cookies.CreateCookies(w, user); err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "error creating cookies", err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
