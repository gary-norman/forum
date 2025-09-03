// Package view contains the HTML rendering logic for the forum application.
package view

import (
	"html/template"
	"path/filepath"

	"github.com/gary-norman/forum/internal/app"
)

type TempHelper struct {
	App *app.App
}

var Template *template.Template

// Init Function to initialise the custom template functions
func (t *TempHelper) Init() {
	tmplFiles1, _ := filepath.Glob("assets/templates/*.html")
	tmplFiles2, _ := filepath.Glob("assets/templates/*.tmpl")
	allFiles := append(tmplFiles1, tmplFiles2...)
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
	}).ParseFiles(allFiles...))
}
