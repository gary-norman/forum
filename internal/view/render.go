// Package view contains the HTML rendering logic for the forum application.
package view

import (
	"html/template"

	"github.com/gary-norman/forum/internal/app"
)

type TempHelper struct {
	App *app.App
}

var Template *template.Template

// init Function to initialise the custom template functions
func (t *TempHelper) Init() {
	Template = template.Must(template.New("").Funcs(template.FuncMap{
		"compareAsInts":  compareAsInts,
		"debugPanic":     debugPanic,
		"decrement":      decrement,
		"dict":           dict,
		"fprint":         fprint,
		"increment":      increment,
		"isValZero":      isValZero,
		"not":            not,
		"or":             or,
		"printType":      printType,
		"random":         randomInt,
		"reactionStatus": t.App.Reactions.GetReactionStatus,
		"same":           checkSameName,
		"startsWith":     startsWith,
	}).ParseGlob("assets/templates/*.html"))
}
