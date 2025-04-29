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
	Close       string
	ConnClose   string
	ConnInit    string
	ConnConn    string
	ConnSuccess string
	DbSuccess   string
	Convert     string
	Cookies     string
	CreateFile  string
	Delete      string
	Divider     string
	Edit        string
	Encode      string
	Execute     string
	Insert      string
	// KeyValuePair: "v" (blue), "v" (white)
	KeyValuePair string
	Login        string
	NoRows       string
	// NotFound: Unable to find "v" called by "v" with error "v"
	NotFound  string
	Open      string
	Parse     string
	Post      string
	Comment   string
	Printf    string
	Protected string
	// Query: Unable to query "v" with error "v"
	Query        string
	Read         string
	RetrieveFile string
	// Register: unable to register "v" with error "v"
	Register  string
	SaveFile  string
	Shutdown  string
	Unmarshal string
	Update    string
	UserModel string
	Write     string
}
type Message struct {
	Name, Text string
}

func JsonError(messageStruct TemplateData) {
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

func JsonPost(messageStruct any) {
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
		DbSuccess:    Colors.Blue + "Database connected: " + Colors.White + "%v " + "v" + Colors.Orange + "%v " + Colors.Green + "- success!\n" + Colors.Reset,
		Convert:      Colors.Red + "Unable to convert " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Cookies:      Colors.Red + "Unable to " + Colors.White + "%v cookies " + Colors.Blue + "with error: " + Colors.Red + "%v" + Colors.Reset,
		CreateFile:   Colors.Red + "Unable to create file " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Delete:       Colors.Red + "Unable to delete " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Divider:      Colors.Grey + "-------------------------------------------------------" + Colors.Reset,
		Edit:         Colors.Red + "Unable to edit " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Encode:       Colors.Red + "Unable to encode JSON called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		Execute:      Colors.Red + "Unable to execute template with error: " + Colors.Red + "%v" + Colors.Reset,
		Insert:       Colors.Red + "Unable to insert " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		KeyValuePair: Colors.Blue + "%v: " + Colors.White + "%v\n",
		Login:        Colors.Red + "Unable to login " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
		NoRows:       Colors.Red + "No rows returned for " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + " with error: " + Colors.Red + "v" + Colors.Reset,
		NotFound: Colors.Red + "Unable to find: " + Colors.White + "%v " + Colors.Blue + "called by " + Colors.
			White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v" + Colors.Reset,
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
