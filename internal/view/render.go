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
		"random":         RandomInt,
		"increment":      Increment,
		"decrement":      Decrement,
		"same":           CheckSameName,
		"compareAsInts":  CompareAsInts,
		"reactionStatus": t.App.Reactions.GetReactionStatus,
		"dict":           dict,
		"isValZero":      isValZero,
		"fprint":         fprint,
		"debugPanic":     debugPanic,
		"or":             or,
		"not":            not,
	}).ParseGlob("assets/templates/*.html"))
}
