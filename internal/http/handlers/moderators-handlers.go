package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gary-norman/forum/internal/app"
	mw "github.com/gary-norman/forum/internal/http/middleware"
)

type ModHandler struct {
	App     *app.App
	Channel *ChannelHandler
	User    *UserHandler
}
type APIResponse struct {
	StatusCode int
	Message    string
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, message string) {
	resp := APIResponse{
		StatusCode: statusCode,
		Message:    message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

func (m *ModHandler) RequestModeration(w http.ResponseWriter, r *http.Request, channelID int64) {
	currentUser, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		log.Printf(ErrorMsgs.NotFound, "currentUser", "requestModeration", "_")
		return
	}
	channelOwner, err := m.App.Channels.GetNameOfChannelOwner(channelID)
	if err != nil {
		log.Printf(ErrorMsgs.KeyValuePair, "error fetching channel owner", err)
	}

	channel, err := m.App.Channels.GetChannelByID(channelID)
	if err != nil {
		log.Printf(ErrorMsgs.KeyValuePair, "error fetching channel", err)
		http.Error(w, `{"error": "channel not found"}`, http.StatusNotFound)
		return
	}

	switch channel.Privacy {
	case true:
		// construct the request, set the status to pending, notify the user
		// send a message to the channel owner
		writeJSONResponse(w, http.StatusOK, fmt.Sprintf("Moderation request sent to %s", channelOwner))
	case false:
		// call the  AddModeration function
		if m.App.Mods.AddModeration(currentUser.ID, channelID) != nil {
			log.Printf(ErrorMsgs.KeyValuePair, "error adding moderation", err)
		}
		writeJSONResponse(w, http.StatusOK, fmt.Sprintf("Welcome to %s!", channel.Name))
	default:
		log.Printf(ErrorMsgs.KeyValuePair, "error determining channel privacy", "switch caught neither true or false")
	}
}
