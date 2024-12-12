package main

import (
	"database/sql"
	"fmt"
	"github.com/gary-norman/forum/internal/sqlite"
	"html/template"
	"log"
	"math/rand"
)

var Tpl *template.Template

// init Function to initialise the custom template functions
func (app *app) init() {

	Tpl = template.Must(template.New("").Funcs(template.FuncMap{
		"random":    RandomInt,
		"increment": Increment,
		"decrement": Decrement,
		"same":      CheckSameName,
		"reactionStatus": func(reactionAuthorID, channelID, passedPostID, passedCommentID int) (sqlite.ReactionStatus, error) {
			// Properly initialise ReactionModel with a valid DB connection
			db, err := sql.Open("sqlite3", "./forum_database.db")
			if err != nil {
				return sqlite.ReactionStatus{}, fmt.Errorf("failed to open database: %w", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(db)

			reactionModel := sqlite.ReactionModel{DB: db}
			return reactionModel.GetReactionStatus(reactionAuthorID, channelID, passedPostID, passedCommentID)
		},

		//Wrap the Method in a Standalone Function
		//Create a wrapper function that calls the method and pass the required arguments
		//without the receiver
		//"isReacted": func(channelID, reactionAuthorID, passedPostID, passedCommentID int) (bool, bool) {
		//	return app.IsReacted(channelID, reactionAuthorID, &passedPostID, &passedCommentID)
		//},
	}).ParseGlob("./assets/templates/*.html"))
}

//// IsReacted Function to check if the user has liked the post or comment, for go templates
//func (app *app) IsReacted(channelID, reactionAuthorID int, passedPostID, passedCommentID *int) (bool, bool) {
//	postID, commentID := 0, 0
//
//	if passedPostID != nil {
//		postID = *passedPostID
//	}
//
//	if passedCommentID != nil {
//		commentID = *passedCommentID
//	}
//
//	// Check if the reaction exists
//	liked, disliked := app.reactions.WhichReaction(reactionAuthorID, channelID, postID, commentID)
//
//	return liked, disliked
//}

// CheckSameName Function to check if the member and artist names are the same, for go templates
func CheckSameName(firstString, secondString string) bool {
	return firstString == secondString
}

// RandomInt Function to get a random integer between 0 and the max number, for go templates
func RandomInt(max int) int {
	return rand.Intn(max)
}

// Increment Function to increment an integer for go templates
func Increment(n int) int {
	return n + 1
}

// Decrement Function to decrement an integer for go templates
func Decrement(n int) int {
	return n - 1
}

// GetTemplate Function to get the template
func GetTemplate() (*template.Template, error) {
	if Tpl == nil {
		return nil, fmt.Errorf("template initialisation failed: template is nil")
	}
	return Tpl, nil
}
