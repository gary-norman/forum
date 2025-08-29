package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gary-norman/forum/internal/app"
	mw "github.com/gary-norman/forum/internal/http/middleware"
	"github.com/gary-norman/forum/internal/models"
	"github.com/gary-norman/forum/internal/view"
)

type ChannelHandler struct {
	App      *app.App
	Reaction *ReactionHandler
	Comment  *CommentHandler
}

func (c *ChannelHandler) GetThisChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(Colors().Orange + "GetThisChannel" + Colors().Reset)
	userLoggedIn := true
	currentUser, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		log.Printf(ErrorMsgs().NotFound, "current user", "getThisChannel", "_")
		userLoggedIn = false
	}

	// Parse channelID from the request
	channelID, err := models.GetIntFromPathValue(r.PathValue("channelId"))
	if err != nil {
		http.Error(w, `{"error": "channelId must be an integer"}`, http.StatusBadRequest)
		return
	}

	// Fetch the channel
	foundChannels, err := c.App.Channels.GetChannelsByID(channelID)
	if err != nil || len(foundChannels) == 0 {
		http.Error(w, `{"error": "Channel not found"}`, http.StatusNotFound)
		return
	}
	thisChannel := foundChannels[0]
	fmt.Printf(ErrorMsgs().KeyValuePair, "Fetching channel", thisChannel.Name)
	models.UpdateTimeSince(&thisChannel)

	// Fetch the channel owner
	thisChannelOwnerName, ownerErr := c.App.Users.GetSingleUserValue(thisChannel.OwnerID, "ID", "username")
	if ownerErr != nil {
		http.Error(w, `{"error": "Error getting channel owner"}`, http.StatusInternalServerError)
		return
	}

	// Fetch channel rules
	thisChannelRules, err := c.App.Rules.AllForChannel(thisChannel.ID)
	if err != nil {
		http.Error(w, `{"error": "Error getting channel rules"}`, http.StatusInternalServerError)
		return
	}

	// Fetch channel posts
	var thisChannelPosts []models.Post
	thisChannelPostIDs, err := c.App.Channels.GetPostIDsFromChannel(thisChannel.ID)
	if err != nil {
		http.Error(w, `{"error": "Error getting Post IDs"}`, http.StatusInternalServerError)
	}
	for p := range thisChannelPostIDs {
		post, err := c.App.Posts.GetPostByID(thisChannelPostIDs[p])
		if err != nil {
			http.Error(w, `{"error": "Error getting post ID:" + thisChannelPostIDs[p]}`, http.StatusInternalServerError)
		}
		thisChannelPosts = append(thisChannelPosts, post)
	}

	allChannels, err := c.App.Channels.All()
	if err != nil {
		http.Error(w, `{"error": "Error getting all channels"}`, http.StatusInternalServerError)
	}
	for c := range allChannels {
		models.UpdateTimeSince(&allChannels[c])
	}
	channelName, err := c.App.Channels.GetChannelNameFromID(thisChannel.ID)
	if err != nil {
		http.Error(w, `{"error": "error fetching channel name"}`, http.StatusInternalServerError)
	}

	// Add channel ID & name and fetch timesince for posts
	for p := range thisChannelPosts {
		thisChannelPosts[p].ChannelID, thisChannelPosts[p].ChannelName = thisChannel.ID, channelName
		// TODO no need for this
		models.UpdateTimeSince(&thisChannelPosts[p])
	}

	// Retrieve total likes and dislikes for each Channel post
	thisChannelPosts = c.Reaction.GetPostsLikesAndDislikes(thisChannelPosts)

	// Retrieve last reaction time for posts
	thisChannelPosts, err = c.Reaction.getLastReactionTimeForPosts(thisChannelPosts)
	if err != nil {
		http.Error(w, `{"error": "Error getting channel posts" }`, http.StatusInternalServerError)
	}

	// Retrieve comments for posts
	thisChannelPosts, err = c.Comment.GetPostsComments(thisChannelPosts)
	if err != nil {
		http.Error(w, `{"error": "Error getting comments" }`, http.StatusInternalServerError)
	}

	ownedChannels := make([]models.Channel, 0)
	joinedChannels := make([]models.Channel, 0)
	ownedAndJoinedChannels := make([]models.Channel, 0)
	isJoined := false
	isOwned := false

	if userLoggedIn {
		isOwned = currentUser.ID == thisChannel.OwnerID
		// attach following/follower numbers to currently logged-in user
		currentUser.Followers, currentUser.Following, err = c.App.Loyalty.CountUsers(currentUser.ID)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > currentUser loyalty", err)
			http.Error(w, `{"error": "Error getting user loyalty"}`, http.StatusInternalServerError)
		}
		// get owned and joined channels of current user
		memberships, memberErr := c.App.Memberships.UserMemberships(currentUser.ID)
		if memberErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "getHome > UserMemberships", memberErr)
			http.Error(w, `{"error": "Error getting user memberships"}`, http.StatusInternalServerError)
		}
		ownedChannels, err = c.App.Channels.OwnedOrJoinedByCurrentUser(currentUser.ID)
		if err != nil {
			log.Printf(ErrorMsgs().Query, "GetThisChannel > user owned channels", err)
			http.Error(w, `{"error": "Error getting user owned channels"}`, http.StatusInternalServerError)
		}
		joinedChannels, err = c.JoinedByCurrentUser(memberships)
		if err != nil {
			log.Printf(ErrorMsgs().Query, "user joined channels", err)
			http.Error(w, `{"error": "Error getting user joined channels"}`, http.StatusInternalServerError)
		}
		ownedAndJoinedChannels = append(ownedChannels, joinedChannels...)

		// Determine whether the current user has joined thisChannel
		for _, channel := range joinedChannels {
			if thisChannel.ID == channel.ID {
				isJoined = true
				break
			}
		}
	}

	data := models.ChannelPage{
		UserID:                 models.NewUUIDField(), // Default value of 0 for logged out users
		CurrentUser:            currentUser,
		Instance:               "channel-page",
		ThisChannel:            thisChannel,
		OwnerName:              thisChannelOwnerName,
		IsOwned:                isOwned,
		IsJoined:               isJoined,
		Rules:                  thisChannelRules,
		Posts:                  thisChannelPosts,
		OwnedChannels:          ownedChannels,
		JoinedChannels:         joinedChannels,
		OwnedAndJoinedChannels: ownedAndJoinedChannels,
		ImagePaths:             c.App.Paths,
	}
	view.RenderPageData(w, data)
}

func (c *ChannelHandler) GetChannelInfoFromPostID(postID int64) (int64, string, error) {
	channelIDs, err := c.App.Channels.GetChannelIdFromPost(postID)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "GetChannelInfoFromPostId > GetChannelIdFromPost", err)
		return 0, "", err
	}
	channelID := channelIDs[0]
	channelName, err := c.App.Channels.GetChannelNameFromID(channelID)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "GetChannelInfoFromPostId > GetChannelNameFromPost", err)
		return 0, "", err
	}
	return channelID, channelName, nil
}

func (c *ChannelHandler) StoreChannel(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
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

	insertErr := c.App.Channels.Insert(
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

func (c *ChannelHandler) StoreMembership(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
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
	joinedChannelID, convErr := strconv.ParseInt(r.PostForm.Get("channelId"), 10, 64)
	if convErr != nil {
		log.Printf(ErrorMsgs().Convert, r.PostForm.Get("channelId"), "StoreMembership > GetChannelID", convErr)
	}
	if err := c.App.Memberships.Insert(user.ID, joinedChannelID); err != nil {
		log.Printf(ErrorMsgs().Post, err)
		http.Error(w, err.Error(), 500)
		return
	}

	channelName, err := c.App.Channels.GetNameOfChannel(joinedChannelID)
	if err != nil {
		log.Printf(ErrorMsgs().Query, "channel", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encErr := json.NewEncoder(w).Encode(map[string]any{
		"code":    http.StatusOK,
		"message": fmt.Sprintf("Welcome to %v!", channelName),
	})
	if encErr != nil {
		log.Printf(ErrorMsgs().Encode, "storeMembership: Accepted", encErr)
		return
	}
}

// JoinedByCurrentUser checks if the currently logged-in user is a member of the current channel
func (c *ChannelHandler) JoinedByCurrentUser(memberships []models.Membership) ([]models.Channel, error) {
	fmt.Println(Colors().Orange + "Checking if this user is a member of this channel" + Colors().Reset)
	fmt.Println(ErrorMsgs().Divider)
	var channels []models.Channel
	for _, membership := range memberships {
		channel, err := c.App.Channels.GetChannelsByID(membership.ChannelID)
		if err != nil {
			return nil, fmt.Errorf(ErrorMsgs().KeyValuePair, "Error calling JoinedByCurrentUser > OwnedOrJoinedByCurrentUser", err)
		}
		channels = append(channels, channel[0])
	}
	// TODO add logic that checks if the user is an owner of this channel
	if len(channels) > 0 {
		fmt.Println(Colors().Green + "Current user is a member of this channel" + Colors().Reset)
	} else {
		fmt.Println(Colors().Red + "Current user is not a member of this channel" + Colors().Reset)
	}
	return channels, nil
}

func (c *ChannelHandler) CreateAndInsertRule(w http.ResponseWriter, r *http.Request) {
	channelID, err := strconv.ParseInt(r.PathValue("channelId"), 10, 64)
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
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "CreateAndInsertRule > id => idInt", err)
		}
		if found {
			err := c.App.Rules.DeleteRule(channelID, idInt)
			if err != nil {
				log.Printf(ErrorMsgs().KeyValuePair, "CreateAndInsertRule > DeleteRule", err)
			}
		} else {
			ruleID, err := c.App.Rules.CreateRule(rule.Rule)
			if err != nil {
				log.Printf(ErrorMsgs().KeyValuePair, "CreateAndInsertRule > CreateRule", err)
			}
			err = c.App.Rules.InsertRule(channelID, ruleID)
			if err != nil {
				log.Printf(ErrorMsgs().KeyValuePair, "CreateAndInsertRule > InsertRule", err)
			}
		}
	}
	http.Redirect(w, r, "/channels/"+r.PathValue("channelId"), http.StatusFound)
}
