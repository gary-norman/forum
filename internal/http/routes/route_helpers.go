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

func NewCommentHandler(app *app.App, reaction *h.ReactionHandler) *h.CommentHandler {
	return &h.CommentHandler{
		App:      app,
		Reaction: reaction,
	}
}

func NewReactionHandler(app *app.App) *h.ReactionHandler {
	return &h.ReactionHandler{
		App: app,
	}
}

func NewChannelHandler(app *app.App, comment *h.CommentHandler, reaction *h.ReactionHandler) *h.ChannelHandler {
	return &h.ChannelHandler{
		App:      app,
		Comment:  comment,
		Reaction: reaction,
	}
}

func NewUserHandler(app *app.App, channel *h.ChannelHandler, comment *h.CommentHandler, reaction *h.ReactionHandler) *h.UserHandler {
	return &h.UserHandler{
		App:      app,
		Comment:  comment,
		Channel:  channel,
		Reaction: reaction,
	}
}

func NewPostHandler(app *app.App, channel *h.ChannelHandler, comment *h.CommentHandler, reaction *h.ReactionHandler) *h.PostHandler {
	return &h.PostHandler{
		App:      app,
		Channel:  channel,
		Comment:  comment,
		Reaction: reaction,
	}
}

func NewHomeHandler(app *app.App, channel *h.ChannelHandler, comment *h.CommentHandler, reaction *h.ReactionHandler) *h.HomeHandler {
	return &h.HomeHandler{
		App:      app,
		Channel:  channel,
		Comment:  comment,
		Reaction: reaction,
	}
}

func NewSessionHandler(app *app.App) *h.SessionHandler {
	return &h.SessionHandler{
		App: app,
	}
}

func NewAuthHandler(app *app.App, session *h.SessionHandler) *h.AuthHandler {
	return &h.AuthHandler{
		App:     app,
		Session: session,
	}
}

func NewSearchHandler(app *app.App) *h.SearchHandler {
	return &h.SearchHandler{
		App: app,
	}
}

func NewRouteHandler(app *app.App) *RouteHandler {
	// Step 1: Create top-level (flat) handlers without nested deps first
	sessionHandler := NewSessionHandler(app)
	reactionHandler := NewReactionHandler(app)
	authHandler := NewAuthHandler(app, sessionHandler)

	// Step 2: Create nested handlers with their deps injected
	commentHandler := NewCommentHandler(app, reactionHandler)
	channelHandler := NewChannelHandler(app, commentHandler, reactionHandler)
	userHandler := NewUserHandler(app, channelHandler, commentHandler, reactionHandler)
	postHandler := NewPostHandler(app, channelHandler, commentHandler, reactionHandler)
	homeHandler := NewHomeHandler(app, channelHandler, commentHandler, reactionHandler)
	searchHandler := NewSearchHandler(app)

	// Step 3: Return fully wired router
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

// func NewRouteHandler(app *app.App) *RouteHandler {
// 	// Step 1: create empty structs
// 	commentHandler := &h.CommentHandler{App: app}
// 	reactionHandler := &h.ReactionHandler{App: app}
// 	channelHandler := &h.ChannelHandler{App: app}
// 	userHandler := &h.UserHandler{App: app}
// 	postHandler := &h.PostHandler{App: app}
// 	homeHandler := &h.HomeHandler{App: app}
// 	sessionHandler := &h.SessionHandler{App: app}
// 	authHandler := &h.AuthHandler{App: app, Session: sessionHandler}
// 	searchHandler := &h.SearchHandler{App: app}
//
// 	// Step 2: wire up shared dependencies
// 	commentHandler.Reaction = reactionHandler
//
// 	channelHandler.Comment = commentHandler
// 	channelHandler.Reaction = reactionHandler
//
// 	userHandler.Channel = channelHandler
// 	userHandler.Comment = commentHandler
// 	userHandler.Reaction = reactionHandler
//
// 	postHandler.Channel = channelHandler
// 	postHandler.Comment = commentHandler
// 	postHandler.Reaction = reactionHandler
//
// 	homeHandler.Channel = channelHandler
// 	homeHandler.Comment = commentHandler
// 	homeHandler.Reaction = reactionHandler
//
// 	return &RouteHandler{
// 		App:      app,
// 		Auth:     authHandler,
// 		Channel:  channelHandler,
// 		Home:     homeHandler,
// 		Post:     postHandler,
// 		Reaction: reactionHandler,
// 		Search:   searchHandler,
// 		Session:  sessionHandler,
// 		User:     userHandler,
// 	}
// }
