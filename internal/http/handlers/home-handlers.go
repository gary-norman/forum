// Package handlers contains the HTTP handlers for the forum application.
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gary-norman/forum/internal/app"
	// "github.com/gary-norman/forum/internal/dao"
	mw "github.com/gary-norman/forum/internal/http/middleware"
	"github.com/gary-norman/forum/internal/models"
	"github.com/gary-norman/forum/internal/view"
)

type HomeHandler struct {
	App      *app.App
	Reaction *ReactionHandler
	Post     *PostHandler
	Comment  *CommentHandler
	Channel  *ChannelHandler
	Mod      *ModHandler
}

var TemplateData models.TemplateData

func (h *HomeHandler) RenderIndex(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	// fmt.Println(r.PathValue("invalidString"))
	// stringCheck := []string{" ", "favicon.ico"}
	// if r.PathValue("invalidString") != "" {
	// 	for _, valid := range stringCheck {
	// 		if r.PathValue("invalidString") != valid {
	// 			view.RenderErrorPage(w, models.NotFoundLocation("home"), 400, models.NotFoundError(r.PathValue("invalidString"), "RenderIndex", fmt.Errorf("the user entered %v after /", r.PathValue("invalidString"))))
	// 			return
	// 		}
	// 	}
	// }

	userLoggedIn := true
	userID := models.ZeroUUIDField()
	// SECTION --- posts and comments ---
	// postDAO := dao.NewDAO[*models.Post](h.App.DB)
	// ctx := context.Background()
	// daoAllPosts, err := postDAO.All(ctx)
	// if err != nil {
	// 	log.Printf(ErrorMsgs.KeyValuePair, "Error fetching daoAllposts", err)
	// } else {
	// 	for _, post := range daoAllPosts {
	// 		fmt.Printf(ErrorMsgs.KeyValuePair, "daoAllPosts", post.Title)
	// 	}
	// }

	// SECTION --- user ---
	allUsers, allUsersErr := h.App.Users.All()
	if allUsersErr != nil {
		log.Printf(ErrorMsgs.Query, "RenderIndex> users > All", allUsersErr)
	}
	for u := range allUsers {
		models.UpdateTimeSince(&allUsers[u])
	}

	// attach following/follower numbers to each user
	for u := range allUsers {
		allUsers[u].Followers, allUsers[u].Following, allUsersErr = h.App.Loyalty.CountUsers(allUsers[u].ID)
		if allUsersErr != nil {
			log.Printf(ErrorMsgs.Query, "RenderIndex> users > All > loyalty", allUsersErr)
		}
	}

	randomUser := GetRandomUser(allUsers)
	currentUser, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		userLoggedIn = false
	}

	var currentUserErr error
	// attach following/follower numbers to the random user
	randomUser.Followers, randomUser.Following, currentUserErr = h.App.Loyalty.CountUsers(randomUser.ID)
	if currentUserErr != nil {
		log.Printf(ErrorMsgs.Query, "RenderIndex> users > All", allUsersErr)
	}

	// SECTION --- posts and comments ---
	allPosts, err := h.App.Posts.All()
	if err != nil {
		log.Printf(ErrorMsgs.KeyValuePair, "Error fetching all posts", err)
	}
	// Retrieve total likes and dislikes for each post
	allPosts = h.Reaction.GetPostsLikesAndDislikes(allPosts)

	// Retrieve last reaction time for posts
	allPosts, err = h.Reaction.getLastReactionTimeForPosts(allPosts)
	if err != nil {
		log.Printf(ErrorMsgs.KeyValuePair, "Error getting last reaction time for allPosts", err)
	}

	for p := range allPosts {
		models.UpdateTimeSince(&allPosts[p])
	}
	allPosts, err = h.Comment.GetPostsComments(allPosts)
	if err != nil {
		log.Printf(ErrorMsgs.NotFound, "allPosts comments", "RenderIndex", err)
	}

	for p := range allPosts {
		channelIDs, err := h.App.Channels.GetChannelIDFromPost(allPosts[p].ID)
		if err != nil {
			log.Printf(ErrorMsgs.KeyValuePair, "RenderIndex > channelID", err)
		}
		if len(allPosts) > 0 && len(channelIDs) > 0 {
			allPosts[p].ChannelID = channelIDs[0]
		} else {
			fetchErr := fmt.Sprintf("post ID: %v does not belong to any channel", allPosts[p].ID)
			fmt.Printf(ErrorMsgs.KeyValuePair, "error fetching posts", fetchErr)
		}
	}

	// SECTION --- channels --
	allChannels, err := h.App.Channels.All()
	if err != nil {
		log.Printf(ErrorMsgs.Query, "channels.All", err)
	}
	for c := range allChannels {
		models.UpdateTimeSince(&allChannels[c])
	}

	for p := range allPosts {
		for _, channel := range allChannels {
			if channel.ID == allPosts[p].ChannelID {
				allPosts[p].ChannelName = channel.Name
			}
		}
	}

	ownedChannels := make([]models.Channel, 0)
	joinedChannels := make([]models.Channel, 0)
	ownedAndJoinedChannels := make([]models.Channel, 0)
	channelMap := make(map[int64]bool)

	if userLoggedIn {
		userID = currentUser.ID
		// attach following/follower numbers to currently logged-in user
		currentUser.Followers, currentUser.Following, err = h.App.Loyalty.CountUsers(currentUser.ID)
		if err != nil {
			fmt.Printf(ErrorMsgs.KeyValuePair, "RenderIndex > currentUser loyalty", err)
		}
		// get owned and joined channels of current user
		memberships, memberErr := h.App.Memberships.UserMemberships(currentUser.ID)
		if memberErr != nil {
			log.Printf(ErrorMsgs.KeyValuePair, "RenderIndex > UserMemberships", memberErr)
		}
		ownedChannels, err = h.App.Channels.OwnedOrJoinedByCurrentUser(currentUser.ID)
		if err != nil {
			log.Printf(ErrorMsgs.Query, "RenderIndex > user owned channels", err)
		}
		joinedChannels, err = h.Channel.JoinedByCurrentUser(memberships)
		if err != nil {
			log.Printf(ErrorMsgs.Query, "user joined channels", err)
		}

		// ownedAndJoinedChannels = append(ownedChannels, joinedChannels...)
		// Add owned channels
		for _, channel := range ownedChannels {
			if !channelMap[channel.ID] {
				channelMap[channel.ID] = true
				ownedAndJoinedChannels = append(ownedAndJoinedChannels, channel)
			}
		}

		// Add joined channels
		for _, channel := range joinedChannels {
			if !channelMap[channel.ID] {
				channelMap[channel.ID] = true
				ownedAndJoinedChannels = append(ownedAndJoinedChannels, channel)
			}
		}
	}

	// SECTION -- template ---
	data := models.TemplateData{
		// ---------- users ----------
		UserID:      userID,
		AllUsers:    allUsers,
		RandomUser:  randomUser,
		CurrentUser: currentUser,
		// ---------- channels ----------
		AllChannels:            allChannels,
		OwnedChannels:          ownedChannels,
		JoinedChannels:         joinedChannels,
		OwnedAndJoinedChannels: ownedAndJoinedChannels,
		// ---------- misc ----------
		Instance:   "home-page",
		ImagePaths: h.App.Paths,
	}
	log.Printf(ErrorMsgs.KeyValuePair, "RenderIndex > data.UserID", data.UserID)
	// models.JsonError(TemplateData)
	tpl, err := view.GetTemplate()
	if err != nil {
		log.Printf(ErrorMsgs.Parse, "templates", "RenderIndex", err)
		return
	}

	t := tpl.Lookup("index.html")

	if t == nil {
		http.Error(w, "Template is not initialized", http.StatusInternalServerError)
		return
	}

	execErr := t.Execute(w, data)
	if execErr != nil {
		log.Printf(ErrorMsgs.Execute, execErr)
		return
	}
	log.Printf(ErrorMsgs.KeyValuePair, "RenderIndex Render time", time.Since(start))
}

func (h *HomeHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	start := time.Now()
	var userPosts []models.Post
	userLoggedIn := true
	userID := models.ZeroUUIDField()

	// SECTION --- user ---
	allUsers, allUsersErr := h.App.Users.All()
	if allUsersErr != nil {
		log.Printf(ErrorMsgs.Query, "getHome> users > All", allUsersErr)
	}
	for u := range allUsers {
		models.UpdateTimeSince(&allUsers[u])
	}

	// attach following/follower numbers to each user
	for u := range allUsers {
		allUsers[u].Followers, allUsers[u].Following, allUsersErr = h.App.Loyalty.CountUsers(allUsers[u].ID)
		if allUsersErr != nil {
			log.Printf(ErrorMsgs.Query, "getHome> users > All > loyalty", allUsersErr)
		}
	}

	randomUser := GetRandomUser(allUsers)
	currentUser, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		userLoggedIn = false
	}

	var currentUserErr error
	// attach following/follower numbers to the random user
	randomUser.Followers, randomUser.Following, currentUserErr = h.App.Loyalty.CountUsers(randomUser.ID)
	if currentUserErr != nil {
		log.Printf(ErrorMsgs.Query, "getHome> users > All", allUsersErr)
	}

	// SECTION --- posts and comments ---
	allPosts, err := h.App.Posts.All()
	if err != nil {
		log.Printf(ErrorMsgs.KeyValuePair, "Error fetching all posts", err)
	}
	// Retrieve total likes and dislikes for each post
	allPosts = h.Reaction.GetPostsLikesAndDislikes(allPosts)

	// Retrieve last reaction time for posts
	allPosts, err = h.Reaction.getLastReactionTimeForPosts(allPosts)
	if err != nil {
		log.Printf(ErrorMsgs.KeyValuePair, "Error getting last reaction time for allPosts", err)
	}

	for p := range allPosts {
		models.UpdateTimeSince(&allPosts[p])
	}
	allPosts, err = h.Comment.GetPostsComments(allPosts)
	if err != nil {
		log.Printf(ErrorMsgs.NotFound, "allPosts comments", "getHome", err)
	}

	for p := range allPosts {
		channelIDs, err := h.App.Channels.GetChannelIDFromPost(allPosts[p].ID)
		if err != nil {
			log.Printf(ErrorMsgs.KeyValuePair, "getHome > channelID", err)
		}
		if len(allPosts) > 0 && len(channelIDs) > 0 {
			allPosts[p].ChannelID = channelIDs[0]
		} else {
			fetchErr := fmt.Sprintf("post ID: %v does not belong to any channel", allPosts[p].ID)
			fmt.Printf(ErrorMsgs.KeyValuePair, "error fetching posts", fetchErr)
		}
	}

	// SECTION --- channels --
	allChannels, err := h.App.Channels.All()
	if err != nil {
		log.Printf(ErrorMsgs.Query, "channels.All", err)
	}
	for c := range allChannels {
		models.UpdateTimeSince(&allChannels[c])
	}

	for p := range allPosts {
		for _, channel := range allChannels {
			if channel.ID == allPosts[p].ChannelID {
				allPosts[p].ChannelName = channel.Name
			}
		}
	}

	ownedChannels := make([]models.Channel, 0)
	joinedChannels := make([]models.Channel, 0)
	ownedAndJoinedChannels := make([]models.Channel, 0)
	channelMap := make(map[int64]bool)

	if userLoggedIn {
		userID = currentUser.ID
		userPosts = h.Post.GetUserPosts(currentUser, allPosts)
		// attach following/follower numbers to currently logged-in user
		currentUser.Followers, currentUser.Following, err = h.App.Loyalty.CountUsers(currentUser.ID)
		if err != nil {
			fmt.Printf(ErrorMsgs.KeyValuePair, "getHome > currentUser loyalty", err)
		}
		// get owned and joined channels of current user
		memberships, memberErr := h.App.Memberships.UserMemberships(currentUser.ID)
		if memberErr != nil {
			log.Printf(ErrorMsgs.KeyValuePair, "getHome > UserMemberships", memberErr)
		}
		ownedChannels, err = h.App.Channels.OwnedOrJoinedByCurrentUser(currentUser.ID)
		if err != nil {
			log.Printf(ErrorMsgs.Query, "GetHome > user owned channels", err)
		}
		joinedChannels, err = h.Channel.JoinedByCurrentUser(memberships)
		if err != nil {
			log.Printf(ErrorMsgs.Query, "user joined channels", err)
		}

		// ownedAndJoinedChannels = append(ownedChannels, joinedChannels...)
		// Add owned channels
		for _, channel := range ownedChannels {
			if !channelMap[channel.ID] {
				channelMap[channel.ID] = true
				ownedAndJoinedChannels = append(ownedAndJoinedChannels, channel)
			}
		}

		// Add joined channels
		for _, channel := range joinedChannels {
			if !channelMap[channel.ID] {
				channelMap[channel.ID] = true
				ownedAndJoinedChannels = append(ownedAndJoinedChannels, channel)
			}
		}

	} else {
		userPosts = allPosts
	}

	// SECTION -- template ---
	data := models.HomePage{
		// ---------- users ----------
		UserID:      userID,
		CurrentUser: currentUser,
		// ---------- posts ----------
		AllPosts:  allPosts,
		UserPosts: userPosts,
		// ---------- channels ----------
		OwnedChannels:          ownedChannels,
		JoinedChannels:         joinedChannels,
		OwnedAndJoinedChannels: ownedAndJoinedChannels,
		// ---------- misc ----------
		Instance:   "home-page",
		ImagePaths: h.App.Paths,
	}

	view.RenderPageData(w, data)
	log.Printf(ErrorMsgs.KeyValuePair, "GetHome Render time", time.Since(start))
}
