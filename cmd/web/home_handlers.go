package main

import (
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"log"
	"net/http"
	"time"
)

// function filters all posts to find posts belonging to the currently logged in user and reners the page with injected home-page template
func (app *app) navigateHome(w http.ResponseWriter, r *http.Request) {
	var userPosts []models.Post

	start := time.Now()
	// SECTION --- posts and comments ---
	// var userLoggedIn bool
	userLoggedIn := true

	allPosts, err := app.posts.All()
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error fetching all posts", err)
	}
	// Retrieve total likes and dislikes for each post
	allPosts = app.getPostsLikesAndDislikes(allPosts)

	for p := range allPosts {
		models.UpdateTimeSince(&allPosts[p])
	}

	allPosts, err = app.getPostsComments(allPosts)
	if err != nil {
		log.Printf(ErrorMsgs().NotFound, "allPosts comments", "getHome", err)
	}

	for p := range allPosts {
		channelIDs, err := app.channels.GetChannelIdFromPost(allPosts[p].ID)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > channelID", err)
		}
		if len(allPosts) > 0 && len(channelIDs) > 0 {
			allPosts[p].ChannelID = channelIDs[0]
		} else {
			fmt.Printf(ErrorMsgs().KeyValuePair, "error fetching posts", "no posts or channel IDs found")
		}
	}

	currentUser, ok := getUserFromContext(r.Context())
	if !ok {
		userLoggedIn = false
	}

	if userLoggedIn {
		userPosts = app.getUserPosts(currentUser, allPosts)
	} else {
		userPosts = allPosts
	}

	data := models.ChannelPage{
		CurrentUser: currentUser,
		Instance:    "home-page",
		Posts:       userPosts,
		ImagePaths:  app.paths,
	}

	renderPageData(w, data)
	log.Printf(ErrorMsgs().KeyValuePair, "Home page rendered in:", time.Since(start))
}

// getUserPosts returns a slice of Posts that belong to channels the user follows. If no user is logged in, it returns all posts
func (app *app) getUserPosts(user *models.User, allPosts []models.Post) []models.Post {
	if user == nil {
		return allPosts
	}
	var postsInUserChannels []models.Post

	userChannels, memberErr := app.memberships.UserMemberships(user.ID)
	if memberErr != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error getting user memberships in response func: ", memberErr)
	}

	// Create a set of user's channel IDs for efficient lookup
	userChannelIDSet := make(map[int]bool)
	for _, membership := range userChannels {
		userChannelIDSet[membership.ChannelID] = true
	}

	// Iterate over all the posts and check if their ChannelID is in the user's channel set
	for _, post := range allPosts {
		if _, exists := userChannelIDSet[post.ChannelID]; exists {
			postsInUserChannels = append(postsInUserChannels, post)
		}
	}

	return postsInUserChannels
}
