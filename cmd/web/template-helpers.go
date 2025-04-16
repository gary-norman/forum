package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"strconv"
)

// FIXME figure out where to put this

var Template *template.Template

// init Function to initialise the custom template functions
func (app *app) init() {
	Template = template.Must(template.New("").Funcs(template.FuncMap{
		"random":         RandomInt,
		"increment":      Increment,
		"decrement":      Decrement,
		"same":           CheckSameName,
		"compareAsInts":  CompareAsInts,
		"reactionStatus": app.reactions.GetReactionStatus,
		"dict":           dict,
		"isValZero":      isValZero,
	}).ParseGlob("assets/templates/*.html"))
}

// CheckSameName Function to check if the member and artist names are the same, for go templates
func CheckSameName(firstString, secondString string) bool {
	return firstString == secondString
}

// CompareAsInts converts both arguments to integers using ConvertToInt and compares them
func CompareAsInts(a, b any) bool {
	intA, errA := ConvertToInt(a)
	intB, errB := ConvertToInt(b)

	if errA != nil || errB != nil {
		log.Printf("error in conversion: %v", errA)
		log.Printf("error in conversion: %v", errB)
		return false // Return false if conversion fails
	}

	return intA == intB
}

// ConvertToInt converts different variable types into an int
func ConvertToInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case *int: // Handle pointer to int
		if v == nil {
			return 0, strconv.ErrSyntax // Handle nil pointers safely
		}
		return *v, nil
	case int64:
		return int(v), nil
	case *int64: // Handle pointer to int64
		if v == nil {
			return 0, strconv.ErrSyntax
		}
		return int(*v), nil
	case float64:
		return int(v), nil
	case *float64: // Handle pointer to float64
		if v == nil {
			return 0, strconv.ErrSyntax
		}
		return int(*v), nil
	case string:
		return strconv.Atoi(v)
	case *string: // Handle pointer to string
		if v == nil {
			return 0, strconv.ErrSyntax
		}
		return strconv.Atoi(*v)
	default:
		return 0, strconv.ErrSyntax
	}
}

// RandomInt Function to get a random integer between 0 and the max number, for go templates
func RandomInt(max int) int {
	return rand.Intn(max)
}

// dict allows 2 parameters to be passed to the {{template}} in the tmpl
func dict(values ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i < len(values); i += 2 {
		key, _ := values[i].(string)
		m[key] = values[i+1]
	}
	return m
}

func isValZero(val string) bool {
	return len(val) == 0
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
