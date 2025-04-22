package handlers

import (
	"github.com/gary-norman/forum/internal/app"
	h "github.com/gary-norman/forum/internal/http/handlers"
)

type RouteHandler struct {
	App      *app.App
	Auth     *h.AuthHandler
	Channel  *h.ChannelHandler
	Comment  *h.CommentHandler
	Home     *h.HomeHandler
	Post     *h.PostHandler
	Reaction *h.ReactionHandler
	Search   *h.SearchHandler
	Session  *h.SessionHandler
	User     *h.UserHandler
}

func NewHandler(app *app.App) *RouteHandler {
	return &RouteHandler{
		App:      app,
		Auth:     &h.AuthHandler{App: app},
		Channel:  &h.ChannelHandler{App: app},
		Comment:  &h.CommentHandler{App: app},
		Home:     &h.HomeHandler{App: app},
		Post:     &h.PostHandler{App: app},
		Reaction: &h.ReactionHandler{App: app},
		Search:   &h.SearchHandler{App: app},
		Session:  &h.SessionHandler{App: app},
		User:     &h.UserHandler{App: app},
	}
}

func NewRouteHandler(app *app.App) *RouteHandler {
	// Step 1: create empty structs
	commentHandler := &h.CommentHandler{App: app}
	reactionHandler := &h.ReactionHandler{App: app}
	channelHandler := &h.ChannelHandler{App: app}
	userHandler := &h.UserHandler{App: app}
	postHandler := &h.PostHandler{App: app}
	homeHandler := &h.HomeHandler{App: app}
	sessionHandler := &h.SessionHandler{App: app}
	authHandler := &h.AuthHandler{App: app, Session: sessionHandler}
	searchHandler := &h.SearchHandler{App: app}

	// Step 2: wire up shared dependencies
	commentHandler.Reaction = reactionHandler

	channelHandler.Comment = commentHandler
	channelHandler.Reaction = reactionHandler

	userHandler.Channel = channelHandler
	userHandler.Comment = commentHandler
	userHandler.Reaction = reactionHandler

	postHandler.Channel = channelHandler
	postHandler.Comment = commentHandler
	postHandler.Reaction = reactionHandler

	homeHandler.Channel = channelHandler
	homeHandler.Comment = commentHandler
	homeHandler.Reaction = reactionHandler

	return &RouteHandler{
		App:      app,
		Auth:     authHandler,
		Channel:  channelHandler,
		Home:     homeHandler,
		Post:     postHandler,
		Reaction: reactionHandler,
		Search:   searchHandler,
		Session:  sessionHandler,
		User:     userHandler,
	}
}
