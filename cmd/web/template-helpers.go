package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gary-norman/forum/internal/models"
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
		"fprint":         fprint,
		"debugPanic":     debugPanic,
		"or":             or,
		"not":            not,
	}).ParseGlob("assets/templates/*.html"))
}

type PageModel interface {
	GetInstance() string
}

// fprint takes a string and an interface and prints them to the console
func fprint(s string, v any) string {
	fmt.Printf(ErrorMsgs().KeyValuePair, s, v)
	return ""
}

// debugPanic takes an interface and returns a string and an error
func debugPanic(v any) (string, error) {
	return "", fmt.Errorf("TEMPLATE PANIC: %#v", v)
}

// or takes two boolean values and returns true if either is true
func or(a, b bool) bool { return a || b }

// not takes a boolean value and returns its negation
func not(a bool) bool { return !a }

func renderPageData[T PageModel](w http.ResponseWriter, data T) {
	var renderedPage bytes.Buffer
	instance := data.GetInstance()
	HTMLstr := models.ToHTMLVar(instance)

	err := Template.ExecuteTemplate(&renderedPage, instance, data)
	if err != nil {
		errorStr := fmt.Sprintf("Error rendering %v: %v", instance, err)
		log.Printf(errorStr)
		http.Error(w, errorStr, http.StatusInternalServerError)
		return
	}

	// Send the pre-rendered HTML as JSON
	response := map[string]string{
		HTMLstr: renderedPage.String(),
	}

	// Write the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "error encoding JSON"}`, http.StatusInternalServerError)
	}
}

// CheckSameName Function to check if  2 strings are identical, for go templates
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
