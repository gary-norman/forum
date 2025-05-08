package view

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/gary-norman/forum/internal/models"
)

func ErrorMsgs() *models.Errors {
	return models.CreateErrorMessages()
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

// CheckSameName Function to check if  2 strings are identical, for go templates
func checkSameName(firstString, secondString string) bool {
	return firstString == secondString
}

func startsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// CompareAsInts converts both arguments to integers using ConvertToInt and compares them
func compareAsInts(a, b any) bool {
	intA, errA := ConvertToInt(a)
	intB, errB := ConvertToInt(b)

	if errA != nil || errB != nil {
		log.Printf("error in conversion: %v", errA)
		log.Printf("error in conversion: %v", errB)
		return false // Return false if conversion fails
	}

	return intA == intB
}

func printType(name, calledBy string, elem any) string {
	str := fmt.Sprintf("Type of %v with value of %v called by %v", name, elem, calledBy)
	fmt.Printf(ErrorMsgs().KeyValuePair, str, reflect.TypeOf(elem))
	return ""
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
func randomInt(max int) int {
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
func increment(n int) int {
	return n + 1
}

// Decrement Function to decrement an integer for go templates
func decrement(n int) int {
	return n - 1
}

// GetTemplate Function to get the template
func GetTemplate() (*template.Template, error) {
	if Template == nil {
		return nil, fmt.Errorf("template initialisation failed: template is nil")
	}
	return Template, nil
}
