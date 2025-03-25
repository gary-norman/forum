package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gary-norman/forum/internal/models"
)

func (app *app) getThisChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var thisChannel models.Channel
	userLoggedIn := true
	currentUser, ok := getUserFromContext(r.Context())
	if !ok {
		fmt.Printf(ErrorMsgs().NotFound, "current user", "getThisChannel", "_")
		userLoggedIn = false
	}
	channelId, err := strconv.Atoi(r.PathValue("channelId"))
	if err != nil {
		fmt.Printf(ErrorMsgs().KeyValuePair, "convert channelId", err)
	}
	foundChannels, err := app.channels.Search("id", channelId)
	if err != nil {
		fmt.Printf(ErrorMsgs().KeyValuePair, "getHome > found channels", err)
	}
	if len(foundChannels) > 0 {
		thisChannel = foundChannels[0]
	} else {
		log.Printf(ErrorMsgs().KeyValuePair, "No channels found", channelId)
	}
	thisChannelOwnerName, ownerErr := app.users.GetSingleUserValue(thisChannel.OwnerID, "ID", "username")
	if ownerErr != nil {
		log.Printf(ErrorMsgs().Query, "getHome > GetSingleUserValue", ownerErr)
	}
	// get all rules for the current channel
	thisChannelRules, err := app.rules.AllForChannel(thisChannel.ID)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "getHome > rules.AllForChannel", err)
	}
	// TODO get channel moderators
	var thisChannelPosts []models.Post
	thisChannelPostIDs, err := app.channels.GetPostIDsFromChannel(thisChannel.ID)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "getHome > channels.GetPostIDsFromChannel", err)
	}
	for _, postID := range thisChannelPostIDs {
		post, err := app.posts.GetPostByID(postID)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > allPosts.GetPostByID", err)
		}
		thisChannelPosts = append(thisChannelPosts, post)
	}
	// Retrieve total likes and dislikes for each Channel post
	thisChannelPosts = app.getPostsLikesAndDislikes(thisChannelPosts)
	thisChannelPosts, err = app.getPostsComments(thisChannelPosts)
	if err != nil {
		log.Printf(ErrorMsgs().NotFound, "thisChannelPosts comments", "getThisChannel", err)
	}
	isOwned := currentUser.ID == thisChannel.OwnerID

	var ownedChannels []models.Channel
	var joinedChannels []models.Channel
	var ownedAndJoinedChannels []models.Channel
	isJoinedOrOwned := false

	if userLoggedIn {
		// attach following/follower numbers to currently logged-in user
		currentUser.Followers, currentUser.Following, err = app.loyalty.CountUsers(currentUser.ID)
		if err != nil {
			fmt.Printf(ErrorMsgs().KeyValuePair, "getHome > currentUser loyalty", err)
		}
		// get owned and joined channels of current user
		memberships, memberErr := app.memberships.UserMemberships(currentUser.ID)
		if memberErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > UserMemberships", memberErr)
		}
		ownedChannels, err = app.channels.OwnedOrJoinedByCurrentUser(currentUser.ID, "OwnerID")
		if err != nil {
			log.Printf(ErrorMsgs().Query, "user owned channels", err)
		}
		joinedChannels, err = app.JoinedByCurrentUser(memberships)
		if err != nil {
			log.Printf(ErrorMsgs().Query, "user joined channels", err)
		}
		ownedAndJoinedChannels = append(ownedChannels, joinedChannels...)
		joined := false
		for _, channel := range joinedChannels {
			if thisChannel.ID == channel.ID {
				joined = true
				break
			}
		}
		isJoinedOrOwned = isOwned || joined
	}

	TemplateData.ThisChannel = thisChannel
	TemplateData.ThisChannelOwnerName = thisChannelOwnerName
	TemplateData.ThisChannelRules = thisChannelRules
	TemplateData.ThisChannelPosts = thisChannelPosts
	TemplateData.ThisChannelIsOwned = isOwned
	TemplateData.OwnedAndJoinedChannels = ownedAndJoinedChannels
	TemplateData.ThisChannelIsOwnedOrJoined = isJoinedOrOwned
}

func (app *app) storeChannel(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
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
	user, ok := getUserFromContext(r.Context())
	if !ok {
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
	encErr := json.NewEncoder(w).Encode(map[string]any{
		"code":    http.StatusOK,
		"message": fmt.Sprintf("Welcome to %v!", channel.Name),
	})
	if encErr != nil {
		log.Printf(ErrorMsgs().Encode, "storeMembership: Accepted", encErr)
		return
	}
}

// func (app *app) requestMembership(w http.ResponseWriter, r *http.Request, userID, channelID int) {
// }

// JoinedByCurrentUser checks if the currently logged-in user is a member of the current channel
func (app *app) JoinedByCurrentUser(memberships []models.Membership) ([]models.Channel, error) {
	fmt.Println(Colors().Orange + "Checking if this user is a member of this channel" + Colors().Reset)
	fmt.Printf(ErrorMsgs().Divider)
	var channels []models.Channel
	for _, membership := range memberships {
		channel, err := app.channels.OwnedOrJoinedByCurrentUser(membership.ChannelID, "ID")
		if err != nil {
			return nil, fmt.Errorf(ErrorMsgs().KeyValuePair, "Error calling JoinedByCurrentUser > OwnedOrJoinedByCurrentUser", err)
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
	var rules []models.PostRule
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
		if found {
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
	http.Redirect(w, r, "/channels/"+r.PathValue("channelId"), http.StatusFound)
}
