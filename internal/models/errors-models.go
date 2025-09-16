package models

import (
	"fmt"
	"reflect"
)

type Colors struct {
	Reset  string
	Red    string
	Green  string
	Yellow string
	Orange string
	Blue   string
	Purple string
	Cyan   string
	Grey   string
	White  string
}

func CreateColors() *Colors {
	colors := &Colors{
		Reset:  "\033[0m",
		Red:    "\033[31m",
		Green:  "\033[32m",
		Yellow: "\033[33m",
		Orange: "\033[38;5;208m",
		Blue:   "\033[34m",
		Purple: "\033[35m",
		Cyan:   "\033[36m",
		Grey:   "\033[37m",
		White:  "\033[97m",
	}
	return colors
}

type Errors struct {
	// Close: Unable to close "v" called by "v" with error "v"
	Close string
	// ConnClose: Unable to close connection to "v" called by "v"
	ConnClose string
	// ConnInit: Unable to initialise connection to "v" called by "v"
	ConnInit string
	// ConnConn: Unable to connect to "v" called by "v(v)"
	ConnConn string
	// ConnSuccess: Server listening on "v" - success!
	ConnSuccess string
	// DBSuccess: Database connected: "v" "v" - success!
	DBSuccess string
	// Convert: Unable to convert "v" called by "v" with error "v"
	Convert string
	// Cookies: Unable to "v" cookies with error "v"
	Cookies string
	// CreateFile: Unable to create file "v" with error "v"
	CreateFile string
	// Delete: Unable to delete "v" with error "v"
	Delete string
	// Divider: Error dividing "v" by "v" with error "v"
	Divider string
	// Edit: Unable to edit "v" with error "v"
	Edit string
	// Encode: Unable to encode "v" with error "v"
	Encode string
	// Execute: Unable to execute "v" with error "v"
	Execute string
	// Generic: Generic error occurred: "v"
	Generic string
	// Insert: Unable to insert "v" with error "v"
	Insert string
	// KeyValuePair: "v" (blue), "v" (white)
	KeyValuePair string
	// Login: Unable to login user "v" with error "v"
	Login string
	// NoRows: No rows found for query "v"
	NoRows string
	// NotFound: Unable to find "v" called by "v" with error "v"
	NotFound string
	// Open: Unable to open "v" with error "v"
	Open string
	// Parse: Unable to parse "v" with error "v"
	Parse string
	// Post: Unable to post "v" with error "v"
	Post string
	// Comment: Unable to add comment "v" with error "v"
	Comment string
	// Printf: Error formatting string "v" with error "v"
	Printf string
	// Protected: Attempt to access protected resource "v" denied
	Protected string
	// Query: Unable to query "v" with error "v"
	Query string
	// Read: Unable to read "v" with error "v"
	Read string
	// RetrieveFile: Unable to retrieve file "v" with error "v"
	RetrieveFile string
	// Register: unable to register "v" with error "v"
	Register string
	// SaveFile: Unable to save file "v" with error "v"
	SaveFile string
	// Shutdown: Error during shutdown: "v"
	Shutdown string
	// Unmarshal: Unable to unmarshal "v" with error "v"
	Unmarshal string
	// Update: Unable to update "v" with error "v"
	Update string
	// UserModel: Error in user model "v" with error "v"
	UserModel string
	// Write: Unable to write "v" with error "v"
	Write string
}

type Message struct {
	Name, Text string
}

func JSONError(messageStruct TemplateData) {
	ErrorMsgs := CreateErrorMessages()
	val := reflect.ValueOf(messageStruct)
	typ := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)
		if fieldType.Name == "CurrentUser" {
			continue
		}
		if fieldType.Name == "Posts" {
			// fmt.Printf(ErrorMsgs.KeyValuePair, "Number of posts", len(field.Interface().([]PostWithWrapping)))
			continue
		}
		if fieldType.Name == "Images" {
			// fmt.Printf(ErrorMsgs.KeyValuePair, "Number of images", len(field.Interface().([]Image)))
			continue
		}
		if fieldType.Name == "Comments" {
			// fmt.Printf(ErrorMsgs.KeyValuePair, "Number of comments", len(field.Interface().([]Comment)))
			continue
		}
		if fieldType.Name == "Reactions" {
			// fmt.Printf(ErrorMsgs.KeyValuePair, "Number of reactions", len(field.Interface().([]Reaction)))
			continue
		}
		fmt.Printf(ErrorMsgs.KeyValuePair, fieldType.Name, field.Interface())
	}
}

func JSONPost(messageStruct any) {
	ErrorMsgs := CreateErrorMessages()
	val := reflect.ValueOf(messageStruct)
	typ := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fmt.Printf(ErrorMsgs.KeyValuePair, fieldType.Name, field.Interface())
	}
}

func CreateErrorMessages() *Errors {
	Colors := CreateColors()
	customErrors := &Errors{
		Close:        Colors.Red + "Unable to close " + Colors.White + "%v" + Colors.Blue + "called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		ConnConn:     Colors.Red + "Unable to connect to " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v(%v)" + Colors.Reset,
		ConnClose:    Colors.Red + "Unable to close connection to " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Reset,
		ConnInit:     Colors.Red + "Unable to initialise connection " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Reset,
		ConnSuccess:  Colors.Blue + "Server listening on " + Colors.White + "%v " + Colors.Green + "- success!\n" + Colors.Reset,
		DBSuccess:    Colors.Blue + "Database connected: " + Colors.White + "%v " + "v" + Colors.Orange + "%v " + Colors.Green + "- success!\n" + Colors.Reset,
		Convert:      Colors.Red + "Unable to convert " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Cookies:      Colors.Red + "Unable to " + Colors.White + "%v cookies " + Colors.Blue + "with error: " + Colors.Red + "%v" + Colors.Reset,
		CreateFile:   Colors.Red + "Unable to create file " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Delete:       Colors.Red + "Unable to delete " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Divider:      Colors.Grey + "-------------------------------------------------------" + Colors.Reset,
		Edit:         Colors.Red + "Unable to edit " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Encode:       Colors.Red + "Unable to encode JSON called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Execute:      Colors.Red + "Unable to execute template with error: " + Colors.Red + "%v" + Colors.Reset,
		Generic:      Colors.Red + "%s: " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Insert:       Colors.Red + "Unable to insert " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		KeyValuePair: Colors.Blue + "%v: " + Colors.White + "%v\n",
		Login:        Colors.Red + "Unable to login " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		NoRows:       Colors.Red + "No rows returned for " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + " with error: " + Colors.Red + "v" + Colors.Reset,
		NotFound:     Colors.Red + "Unable to find: " + Colors.White + "%v " + Colors.Blue + "called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Open:         Colors.Red + "Unable to open " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error " + Colors.Red + "%v" + Colors.Reset,
		Parse:        Colors.Red + "Unable to parse " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Post:         Colors.Red + "Unable to post with error: " + Colors.Red + "%v" + Colors.Reset,
		Comment:      Colors.Red + "Unable to Comment with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Printf:       Colors.Red + "Unable to print with error: " + Colors.Red + "%v" + Colors.Reset,
		Protected:    Colors.Red + "CSRF validation failed for user " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Query:        Colors.Red + "Unable to query " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Read:         Colors.Red + "Unable to read " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		RetrieveFile: Colors.Red + "Unable to retrieve file " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Register:     Colors.Red + "Unable to register with error: " + Colors.Red + "%v" + Colors.Reset,
		SaveFile:     Colors.Red + "Unable to save file " + Colors.White + "%v" + Colors.Blue + " to " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Shutdown:     Colors.Red + "HTTP shutdown error: " + Colors.White + "%v" + Colors.Reset,
		Update:       Colors.Red + "Unable to update " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		UserModel:    Colors.Red + "Usermodel or DB called in " + Colors.White + "%v" + Colors.Blue + " for " + Colors.White + "%v" + Colors.Blue + " is nil" + Colors.Reset,
		Write:        Colors.Red + "Unable to write to " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Reset,
	}
	return customErrors
}
