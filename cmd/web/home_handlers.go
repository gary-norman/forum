package main

import (
	"github.com/gary-norman/forum/internal/models"
	"log"
)

// getUserPosts returns a slice of Posts that belong to channels the user follows. If no user is logged in, it returns all posts
func (app *app) getUserPosts(user *models.User, allPosts []models.Post) []models.Post {
	if user == nil {
		return allPosts
	}
	var ownedChannels, joinedChannels, ownedAndJoinedChannels []models.Channel
	var postsInUserChannels []models.Post

	memberships, memberErr := app.memberships.UserMemberships(user.ID)
	if memberErr != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "getHome > UserMemberships", memberErr)
	}
	ownedChannels, ownedErr := app.channels.OwnedOrJoinedByCurrentUser(user.ID, "OwnerID")
	if ownedErr != nil {
		log.Printf(ErrorMsgs().Query, "user owned channels", ownedErr)
	}
	joinedChannels, joinedErr := app.JoinedByCurrentUser(memberships)
	if joinedErr != nil {
		log.Printf(ErrorMsgs().Query, "user joined channels", joinedErr)
	}
	ownedAndJoinedChannels = append(ownedChannels, joinedChannels...)

	if len(ownedAndJoinedChannels) == 0 {
		return allPosts
	}

	// Create a set of user's channel IDs for efficient lookup
	userChannelIDSet := make(map[int]bool)
	for _, membership := range ownedAndJoinedChannels {
		userChannelIDSet[membership.ID] = true
	}

	// Iterate over all the posts and check if their ChannelID is in the user's channel set
	for _, post := range allPosts {
		if _, exists := userChannelIDSet[post.ChannelID]; exists {
			postsInUserChannels = append(postsInUserChannels, post)
		}
	}

	return postsInUserChannels
}
