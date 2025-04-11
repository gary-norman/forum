package main

import (
	"bytes"
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

	userLoggedIn := true
	currentUser, ok := getUserFromContext(r.Context())
	if !ok {
		log.Printf(ErrorMsgs().NotFound, "current user", "getThisChannel", "_")
		userLoggedIn = false
	}

	// Parse channelId from the request
	channelId, err := models.GetIntFromPathValue(r.PathValue("channelId"))
	if err != nil {
		http.Error(w, `{"error": "channelId must be an integer"}`, http.StatusBadRequest)
		return
	}

	// Fetch the channel
	foundChannels, err := app.channels.SearchChannelsByColumn("ID", channelId)
	if err != nil || len(foundChannels) == 0 {
		http.Error(w, `{"error": "Channel not found"}`, http.StatusNotFound)
		return
	}
	thisChannel := foundChannels[0]
	fmt.Printf(ErrorMsgs().KeyValuePair, "Fetching channel", thisChannel.Name)

	// Fetch the channel owner
	thisChannelOwnerName, ownerErr := app.users.GetSingleUserValue(thisChannel.OwnerID, "ID", "username")
	if ownerErr != nil {
		http.Error(w, `{"error": "Error getting channel owner"}`, http.StatusInternalServerError)
		return
	}

	// Fetch channel rules
	thisChannelRules, err := app.rules.AllForChannel(thisChannel.ID)
	if err != nil {
		http.Error(w, `{"error": "Error getting channel rules"}`, http.StatusInternalServerError)
		return
	}

	// Fetch channel posts
	thisChannelPosts := []models.Post{}
	thisChannelPostIDs, err := app.channels.GetPostIDsFromChannel(thisChannel.ID)
	if err == nil {
		for _, postID := range thisChannelPostIDs {
			post, err := app.posts.GetPostByID(postID)
			if err == nil {
				thisChannelPosts = append(thisChannelPosts, post)
			}
		}
	}

	for p := range thisChannelPosts {
		models.UpdateTimeSince(&thisChannelPosts[p])
	}

	allChannels, err := app.channels.All()
	if err != nil {
		http.Error(w, `{"error": "Error getting all channels"}`, http.StatusInternalServerError)
	}
	for c := range allChannels {
		models.UpdateTimeSince(&allChannels[c])
	}

	for p := range thisChannelPosts {
		channelIDs, err := app.channels.GetChannelIdFromPost(thisChannelPosts[p].ID)
		if err != nil {
			http.Error(w, `{"error": "Error getting channel ID from post"}`, http.StatusInternalServerError)
		}
		thisChannelPosts[p].ChannelID = channelIDs[0]
	}

	for p := range thisChannelPosts {
		for _, channel := range allChannels {
			if channel.ID == thisChannelPosts[p].ChannelID {
				thisChannelPosts[p].ChannelName = channel.Name
			}
		}
	}
	// Retrieve total likes and dislikes for each Channel post
	thisChannelPosts = app.getPostsLikesAndDislikes(thisChannelPosts)
	thisChannelPosts, err = app.getPostsComments(thisChannelPosts)
	if err != nil {
		http.Error(w, `{"error": "Error getting comments" }`, http.StatusInternalServerError)
	}

	var ownedChannels, joinedChannels, ownedAndJoinedChannels []models.Channel
	isJoined := false
	isOwned := false
	isJoinedOrOwned := false

	if userLoggedIn {
		isOwned = currentUser.ID == thisChannel.OwnerID
		// attach following/follower numbers to currently logged-in user
		currentUser.Followers, currentUser.Following, err = app.loyalty.CountUsers(currentUser.ID)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > currentUser loyalty", err)
			http.Error(w, `{"error": "Error getting user loyalty"}`, http.StatusInternalServerError)
		}
		// get owned and joined channels of current user
		memberships, memberErr := app.memberships.UserMemberships(currentUser.ID)
		if memberErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > UserMemberships", memberErr)
			http.Error(w, `{"error": "Error getting user memberships"}`, http.StatusInternalServerError)
		}
		ownedChannels, err = app.channels.OwnedOrJoinedByCurrentUser(currentUser.ID, "OwnerID")
		if err != nil {
			log.Printf(ErrorMsgs().Query, "user owned channels", err)
			http.Error(w, `{"error": "Error getting user owned channels"}`, http.StatusInternalServerError)
		}
		joinedChannels, err = app.JoinedByCurrentUser(memberships)
		if err != nil {
			log.Printf(ErrorMsgs().Query, "user joined channels", err)
			http.Error(w, `{"error": "Error getting user joined channels"}`, http.StatusInternalServerError)
		}
		ownedAndJoinedChannels = append(ownedChannels, joinedChannels...)
		for _, channel := range joinedChannels {
			if thisChannel.ID == channel.ID {
				isJoined = true
				break
			}
		}
		isJoinedOrOwned = isOwned || isJoined
	}

	paths := models.ImagePaths{
		Channel: "db/userdata/images/channel-images",
		Post:    "db/userdata/images/post-images",
		User:    "db/userdata/images/user-images",
	}

	data := models.ChannelPage{
		TestString:             "This is a test string",
		CurrentUser:            currentUser,
		ThisChannel:            thisChannel,
		ThisChannelOwnerName:   thisChannelOwnerName,
		ThisChannelRules:       thisChannelRules,
		ThisChannelPosts:       thisChannelPosts,
		ThisChannelIsOwned:     isOwned,
		OwnedAndJoinedChannels: ownedAndJoinedChannels,
		IsJoinedOrOwned:        isJoinedOrOwned,
		IsPostPage:             false,
		Instance:               "channel-page",
		ImagePaths:             paths,
	}

	fmt.Printf(ErrorMsgs().KeyValuePair, "Channel page avatar", thisChannel.Avatar)
	fmt.Printf(ErrorMsgs().KeyValuePair, "Channel page posts", len(thisChannelPosts))
	fmt.Printf(ErrorMsgs().KeyValuePair, "Paths", paths)

	// data := models.ChannelPageBanner{
	// 	TestString:             "This is a test string",
	// 	ThisChannel:            thisChannel,
	// 	OwnedAndJoinedChannels: ownedAndJoinedChannels,
	// 	IsJoinedOrOwned:        isJoinedOrOwned,
	// 	ThisChannelOwnerName:   thisChannelOwnerName,
	// 	ThisChannelRules:       thisChannelRules,
	// }

	// data := models.Postplus{
	//   ID: post.ID,
	//   Title: post.Title,
	//   Content: post.Content,
	//   Images: post.Images,
	//   Created: post.Created,
	//   TimeSince: post.TimeSince,
	//   IsCommentable: post.IsCommentable,
	//   Author: post.Author,
	//   AuthorID: post.AuthorID,
	//   AuthorAvatar: post.AuthorAvatar,
	//   ChannelID: thisChannel.ID,
	//   ChannelName: thisChannel.Name,
	//   IsFlagged: post.IsFlagged,
	//   Likes: post.Likes,
	//   Dislikes: post.Dislikes,
	//   CommentsCount: post.CommentsCount,
	//   Comments: post.Comments,
	//   IsPostPage: false,
	//   Instance: "channel-page",
	// }

	// Render the `post-card.html` subtemplate
	var renderedChannelPage bytes.Buffer
	postsErr := Template.ExecuteTemplate(&renderedChannelPage, "channel-page", data)
	if postsErr != nil {
		http.Error(w, "Error rendering channel-page", http.StatusInternalServerError)
		return
	}

	// Send the pre-rendered HTML as JSON
	response := map[string]string{
		"postsHTML": renderedChannelPage.String(),
	}
	log.Printf(ErrorMsgs().KeyValuePair, len(thisChannelPosts), thisChannel.Name)
	json.NewEncoder(w).Encode(response)
}

func (app *app) GetChannelInfoFromPostID(postID int) (int, string, error) {
	channelIDs, err := app.channels.GetChannelIdFromPost(postID)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "GetChannelInfoFromPostId > GetChannelIdFromPost", err)
		return 0, "", err
	}
	channelID := channelIDs[0]
	channelName, err := app.channels.GetChannelNameFromID(channelID)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "GetChannelInfoFromPostId > GetChannelNameFromPost", err)
		return 0, "", err
	}
	return channelID, channelName, nil
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
	channels, err := app.channels.SearchChannelsByColumn("id", joinedChannelID)
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
