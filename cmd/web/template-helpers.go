package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
)

var Tpl *template.Template

// init Function to initialise the custom template functions
func init() {
	IsLiked := (*app).IsLiked
	Tpl = template.Must(template.New("").Funcs(template.FuncMap{
		"random":    RandomInt,
		"increment": Increment,
		"decrement": Decrement,
		"same":      CheckSameName,
		"isLiked":   IsLiked,
	}).ParseGlob("./assets/templates/*.html"))
}

// IsLiked Function to check if the user has liked the post, for go templates
func (app *app) IsLiked(channelID, authorID, userID int, passedPostID, passedCommentID *int) (bool, error) {
	postID, commentID := 0, 0

	if passedPostID != nil {
		//log.Println("ReactedPostID:", *reactionData.ReactedPostID)
		postID = *passedPostID
	}

	if passedCommentID != nil {
		//log.Printf("ReactedCommentID: %d", *reactionData.ReactedPostID)
		commentID = *passedCommentID
	}

	// Check if the user already reacted (like/dislike) and update or delete the reaction if needed
	existingReaction, err := app.reactions.CheckExistingReaction(authorID, channelID, postID, commentID)
	if err != nil {
		// Use your custom error message for fetching errors
		log.Printf("Error fetching existing reaction: %v", err)
		return false, err
	}

	if existingReaction != nil {
		log.Printf("Existing Reaction: %+v", existingReaction)
		return true, nil
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
