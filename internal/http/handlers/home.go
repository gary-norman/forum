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

func (h *HomeHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var userPosts []models.Post
	userLoggedIn := true
	// SECTION --- posts and comments ---
	// postDAO := dao.NewDAO[*models.Post](h.App.Db)
	// ctx := context.Background()
	// daoAllPosts, err := postDAO.All(ctx)
	// if err != nil {
	// 	log.Printf(ErrorMsgs().KeyValuePair, "Error fetching daoAllosts", err)
	// } else {
	// 	for _, post := range daoAllPosts {
	// 		fmt.Printf(ErrorMsgs().KeyValuePair, "daoAllPosts", post.Title)
	// 	}
	// }

	// SECTION --- user ---
	allUsers, allUsersErr := h.App.Users.All()
	if allUsersErr != nil {
		log.Printf(ErrorMsgs().Query, "getHome> users > All", allUsersErr)
	}
	for u := range allUsers {
		models.UpdateTimeSince(&allUsers[u])
	}

	// attach following/follower numbers to each user
	for u := range allUsers {
		allUsers[u].Followers, allUsers[u].Following, allUsersErr = h.App.Loyalty.CountUsers(allUsers[u].ID)
		if allUsersErr != nil {
			log.Printf(ErrorMsgs().Query, "getHome> users > All > loyalty", allUsersErr)
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
		log.Printf(ErrorMsgs().Query, "getHome> users > All", allUsersErr)
	}

	// SECTION --- posts and comments ---
	allPosts, err := h.App.Posts.All()
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error fetching all posts", err)
	}
	// Retrieve total likes and dislikes for each post
	allPosts = h.Reaction.GetPostsLikesAndDislikes(allPosts)

	for p := range allPosts {
		models.UpdateTimeSince(&allPosts[p])
	}
	allPosts, err = h.Comment.GetPostsComments(allPosts)
	if err != nil {
		log.Printf(ErrorMsgs().NotFound, "allPosts comments", "getHome", err)
	}

	for p := range allPosts {
		channelIDs, err := h.App.Channels.GetChannelIdFromPost(allPosts[p].ID)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > channelID", err)
		}
		if len(allPosts) > 0 && len(channelIDs) > 0 {
			allPosts[p].ChannelID = channelIDs[0]
		} else {
			fmt.Printf(ErrorMsgs().KeyValuePair, "error fetching posts", "no posts or channel IDs found")
		}
	}

	// SECTION --- channels --
	allChannels, err := h.App.Channels.All()
	if err != nil {
		log.Printf(ErrorMsgs().Query, "channels.All", err)
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

	var ownedChannels, joinedChannels, ownedAndJoinedChannels []models.Channel
	if userLoggedIn {
		userPosts = h.Post.GetUserPosts(currentUser, allPosts)
		// attach following/follower numbers to currently logged-in user
		currentUser.Followers, currentUser.Following, err = h.App.Loyalty.CountUsers(currentUser.ID)
		if err != nil {
			fmt.Printf(ErrorMsgs().KeyValuePair, "getHome > currentUser loyalty", err)
		}
		// get owned and joined channels of current user
		memberships, memberErr := h.App.Memberships.UserMemberships(currentUser.ID)
		if memberErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > UserMemberships", memberErr)
		}
		ownedChannels, err = h.App.Channels.OwnedOrJoinedByCurrentUser(currentUser.ID)
		if err != nil {
			log.Printf(ErrorMsgs().Query, "user owned channels", err)
		}
		joinedChannels, err = h.Channel.JoinedByCurrentUser(memberships)
		if err != nil {
			log.Printf(ErrorMsgs().Query, "user joined channels", err)
		}
		ownedAndJoinedChannels = append(ownedChannels, joinedChannels...)
	} else {
		userPosts = allPosts
	}

	// SECTION -- template ---
	TemplateData := models.TemplateData{
		// ---------- users ----------
		UserID:      models.NewUUIDField(), // Default value of 0 for logged out users
		AllUsers:    allUsers,
		RandomUser:  randomUser,
		CurrentUser: currentUser,
		// ---------- posts ----------
		Posts:     allPosts,
		UserPosts: userPosts,
		// ---------- channels ----------
		AllChannels:            allChannels,
		OwnedChannels:          ownedChannels,
		JoinedChannels:         joinedChannels,
		OwnedAndJoinedChannels: ownedAndJoinedChannels,
		// ---------- misc ----------
		Instance:   "home-page",
		ImagePaths: h.App.Paths,
	}
	// models.JsonError(TemplateData)
	tpl, err := view.GetTemplate()
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
	log.Printf(ErrorMsgs().KeyValuePair, "GetHome Render time:", time.Since(start))
}
