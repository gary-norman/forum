package main

import (
	"fmt"
	"html/template"
	"math/rand"
)

var Template *template.Template

// init Function to initialise the custom template functions
func (app *app) init() {

	Template = template.Must(template.New("").Funcs(template.FuncMap{
		"random":         RandomInt,
		"increment":      Increment,
		"decrement":      Decrement,
		"same":           CheckSameName,
		"reactionStatus": app.reactions.GetReactionStatus,
	}).ParseGlob("./assets/templates/*.html"))
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
	if Template == nil {
		return nil, fmt.Errorf("template initialisation failed: template is nil")
	}
	return Template, nil
}
