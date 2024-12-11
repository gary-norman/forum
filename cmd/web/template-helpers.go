package main

import (
	"fmt"
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
		//Wrap the Method in a Standalone Function
		//Create a wrapper function that calls the method and pass the required arguments
		//without the receiver
		"isReacted": func(channelID, reactionAuthorID, passedPostID, passedCommentID int) (bool, error) {
			return app.IsReacted(channelID, reactionAuthorID, &passedPostID, &passedCommentID)
		},
	}).ParseGlob("./assets/templates/*.html"))
}

// IsReacted Function to check if the user has liked the post or comment, for go templates
func (app *app) IsReacted(channelID, reactionAuthorID int, passedPostID, passedCommentID *int) (bool, error) {
	postID, commentID := 0, 0

	if passedPostID != nil {
		postID = *passedPostID
	}

	if passedCommentID != nil {
		commentID = *passedCommentID
	}

	// Check if the reaction exists
	exists, err := app.reactions.Exists(reactionAuthorID, channelID, postID, commentID)
	if err != nil {
		// Use your custom error message for fetching errors
		log.Printf("Error fetching existing reaction: %v", err)
		return false, err
	}
	if exists {
		return exists, nil
	}

	return false, nil
}

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
