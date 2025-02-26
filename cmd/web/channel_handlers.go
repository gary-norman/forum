package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gary-norman/forum/cmd/web"
	"github.com/gary-norman/forum/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (app *main.app) getThisChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var thisChannel models.ChannelWithDaysAgo
	channelId, err := strconv.Atoi(r.PathValue("channelId"))
	if err != nil {
		fmt.Printf(main.ErrorMsgs().KeyValuePair, "convert channelId", err)
	}
	foundChannels, err := app.channels.Search("id", channelId)
	if err != nil {
		fmt.Printf(main.ErrorMsgs().KeyValuePair, "getHome > found channels", err)
	}
	if len(foundChannels) > 0 {
		thisChannel = models.ChannelWithDaysAgo{
			Channel:   foundChannels[0],
			TimeSince: main.getTimeSince(foundChannels[0].Created),
		}
	}
	thisChannelOwnerName, ownerErr := app.users.GetSingleUserValue(thisChannel.OwnerID, "ID", "username")
	if ownerErr != nil {
		log.Printf(main.ErrorMsgs().Query, "getHome > GetSingleUserValue", ownerErr)
	}
	main.TemplateData.ThisChannel = thisChannel
	main.TemplateData.ThisChannelOwnerName = thisChannelOwnerName
}

func (app *main.app) storeChannel(w http.ResponseWriter, r *http.Request) {
	user, getUserErr := app.GetLoggedInUser(r)
	if getUserErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	parseErr := r.ParseMultipartForm(10 << 20)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), 400)
		log.Printf(main.ErrorMsgs().Parse, "storeChannel", parseErr)
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
	createChannelData.Avatar = main.GetFileName(r, "file-drop", "storeChannel", "channel")

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
		log.Printf(main.ErrorMsgs().Post, insertErr)
		http.Error(w, insertErr.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *main.app) storeMembership(w http.ResponseWriter, r *http.Request) {
	user, getUserErr := app.GetLoggedInUser(r)
	if getUserErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if parseErr := r.ParseForm(); parseErr != nil {
		http.Error(w, parseErr.Error(), 400)
		log.Printf(main.ErrorMsgs().Parse, "storeMembership", parseErr)
		return
	}
	fmt.Printf("user: %v", user.Username)
	// get channelID
	joinedChannelID, convErr := strconv.Atoi(r.PostForm.Get("channelId"))
	if convErr != nil {
		log.Printf(main.ErrorMsgs().Convert, r.PostForm.Get("channelId"), "StoreMembership > GetChannelID", convErr)
		log.Printf("Unable to convert %v to integer\n", r.PostForm.Get("channelId"))
	}
	// get slice of channels (in this case it is only 1, but the function still returns a slice)
	channels, err := app.channels.Search("id", joinedChannelID)
	if err != nil {
		log.Printf(main.ErrorMsgs().Query, "channel", err)
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
		log.Printf(main.ErrorMsgs().Post, insertErr)
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
		log.Printf(main.ErrorMsgs().Encode, "storeMembership: Accepted", encErr)
		return
	}
}

// func (app *app) requestMembership(w http.ResponseWriter, r *http.Request, userID, channelID int) {
// }

// JoinedByCurrentUser checks if the currently logged-in user is a member of the current channel
func (app *main.app) JoinedByCurrentUser(memberships []models.Membership) ([]models.Channel, error) {
	fmt.Println(main.Colors().Orange + "Checking if this user is a member of this channel" + main.Colors().Reset)
	fmt.Printf(main.ErrorMsgs().Divider)
	var channels []models.Channel
	for _, membership := range memberships {
		channel, err := app.channels.OwnedOrJoinedByCurrentUser(membership.ChannelID, "ID")
		if err != nil {
			return nil, fmt.Errorf(main.ErrorMsgs().KeyValuePair, "Error calling JoinedByCurrentUser > OwnedOrJoinedByCurrentUser", err)
		}
		channels = append(channels, channel[0])
	}
	if len(channels) > 0 {
		fmt.Println(main.Colors().Green + "Current user is a member of this channel" + main.Colors().Reset)
	} else {
		fmt.Println(main.Colors().Red + "Current user is not a member of this channel" + main.Colors().Reset)
	}
	return channels, nil
}

func (app *main.app) CreateAndInsertRule(w http.ResponseWriter, r *http.Request) {
	channelId, err := strconv.Atoi(r.PathValue("channelId"))
	if err != nil {
		log.Printf(main.ErrorMsgs().KeyValuePair, "CreateAndInsertRule > convert channelId to int", err)
	}

	// Get the "rules" input value
	rulesJSON := r.FormValue("rules")
	if rulesJSON == "" { // TODO send this message to the user
		log.Printf(main.ErrorMsgs().KeyValuePair, "message to user", "you have not added or removed any rules")
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
			log.Printf(main.ErrorMsgs().KeyValuePair, "CreateAndInsertRule > id => idInt", err)
		}
		if found {
			err := app.rules.DeleteRule(channelId, idInt)
			if err != nil {
				log.Printf(main.ErrorMsgs().KeyValuePair, "CreateAndInsertRule > DeleteRule", err)
			}
		} else {
			ruleId, err := app.rules.CreateRule(rule.Text)
			if err != nil {
				log.Printf(main.ErrorMsgs().KeyValuePair, "CreateAndInsertRule > CreateRule", err)
			}
			err = app.rules.InsertRule(channelId, ruleId)
			if err != nil {
				log.Printf(main.ErrorMsgs().KeyValuePair, "CreateAndInsertRule > InsertRule", err)
			}
		}
	}
	http.Redirect(w, r, "/channels/"+r.PathValue("channelId"), http.StatusFound)
}
