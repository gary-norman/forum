package main

import (
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"log"
	"net/http"
	"strconv"
)

func (app *app) getThisChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var thisChannel models.ChannelWithDaysAgo
	channelId, err := strconv.Atoi(r.PathValue("channelId"))
	if err != nil {
		fmt.Printf(ErrorMsgs().KeyValuePair, "convert channelId", err)
	}
	foundChannels, err := app.channels.Search("id", channelId)
	if err != nil {
		fmt.Printf(ErrorMsgs().KeyValuePair, "getHome > found channels", err)
	}
	if len(foundChannels) > 0 {
		thisChannel = models.ChannelWithDaysAgo{
			Channel:   foundChannels[0],
			TimeSince: getTimeSince(foundChannels[0].Created),
		}
	}
	thisChannelOwnerName, ownerErr := app.users.GetSingleUserValue(thisChannel.OwnerID, "ID", "username")
	if ownerErr != nil {
		log.Printf(ErrorMsgs().Query, "getHome > GetSingleUserValue", ownerErr)
	}

	TemplateData.ThisChannel = thisChannel
	TemplateData.ThisChannelOwnerName = thisChannelOwnerName

}
