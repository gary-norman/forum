package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (app *app) search(w http.ResponseWriter, r *http.Request) {
	//start := time.Now()

	// SECTION --- channels --
	allChannels, err := app.channels.All()
	if err != nil {
		log.Printf(ErrorMsgs().Query, "channels.All", err)
	}

	// SECTION --- posts ---
	// var userLoggedIn bool
	allPosts, err := app.posts.All()
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error fetching all posts", err)
	}

	for p := range allPosts {
		channelIDs, err := app.channels.GetChannelIdFromPost(allPosts[p].ID)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > channelID", err)
		}
		allPosts[p].ChannelID = channelIDs[0]
	}

	for p := range allPosts {
		for _, channel := range allChannels {
			if channel.ID == allPosts[p].ChannelID {
				allPosts[p].ChannelName = channel.Name
			}
		}
	}

	// SECTION --- user ---
	allUsers, allUsersErr := app.users.All()
	if allUsersErr != nil {
		log.Printf(ErrorMsgs().Query, "getHome> users > All", allUsersErr)
	}

	searchResults := map[string]interface{}{
		"users":    allUsers,
		"channels": allChannels,
		"posts":    allPosts,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(searchResults); err != nil {
		log.Printf(ErrorMsgs().Encode, "search results", err)
		http.Error(w, "Error encoding search results", http.StatusInternalServerError)
		return
	}

	//log.Printf(ErrorMsgs().KeyValuePair, "Search data fetched in:", time.Since(start))
}
