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
	Close        string
	ConnClose    string
	ConnInit     string
	ConnConn     string
	Cookies      string
	CreateFile   string
	Divider      string
	Edit         string
	Execute      string
	KeyValuePair string
	Login        string
	NoRows       string
	NotFound     string
	Open         string
	Parse        string
	Post         string
	Printf       string
	Protected    string
	Query        string
	Read         string
	RetrieveFile string
	Register     string
	SaveFile     string
	Unmarshal    string
	UserModel    string
	Write        string
}
type Message struct {
	Name, Text string
}

func JsonError(messageStruct TemplateData) {
	ErrorMsgs := CreateErrorMessages()
	val := reflect.ValueOf(messageStruct)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		if fieldType.Name == "CurrentUser" {
			continue
		}
		if fieldType.Name == "Posts" {
			fmt.Printf(ErrorMsgs.KeyValuePair, "Number of posts", len(field.Interface().([]PostWithDaysAgo)))
			continue
		}
		if fieldType.Name == "Images" {
			fmt.Printf(ErrorMsgs.KeyValuePair, "Number of images", len(field.Interface().([]Image)))
			continue
		}
		if fieldType.Name == "Comments" {
			fmt.Printf(ErrorMsgs.KeyValuePair, "Number of comments", len(field.Interface().([]Comment)))
			continue
		}
		if fieldType.Name == "Reactions" {
			fmt.Printf(ErrorMsgs.KeyValuePair, "Number of reactions", len(field.Interface().([]Reaction)))
			continue
		}
		if fieldType.Name == "NotifyPlaceHolder" {
			JsonNotifyPlaceholder(field.Interface().(NotifyPlaceholder))
			continue
		}
		fmt.Printf(ErrorMsgs.KeyValuePair, fieldType.Name, field.Interface())
	}
}
func JsonPost(messageStruct PostWithDaysAgo) {
	ErrorMsgs := CreateErrorMessages()
	val := reflect.ValueOf(messageStruct)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fmt.Printf(ErrorMsgs.KeyValuePair, fieldType.Name, field.Interface())
	}
}
func JsonNotifyPlaceholder(messageStruct NotifyPlaceholder) {
	ErrorMsgs := CreateErrorMessages()
	val := reflect.ValueOf(messageStruct)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fmt.Printf(ErrorMsgs.KeyValuePair, fieldType.Name, field.Interface())
	}
}

func CreateErrorMessages() *Errors {
	Colors := CreateColors()
	customErrors := &Errors{
		Close:        Colors.Blue + "Unable to close " + Colors.White + "%v" + Colors.Blue + "called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		ConnConn:     Colors.Blue + "Unable to connect to " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v\n" + Colors.Reset,
		ConnClose:    Colors.Blue + "Unable to close connection to " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v\n" + Colors.Reset,
		ConnInit:     Colors.Blue + "Unable to initialise connection " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v\n" + Colors.Reset,
		Cookies:      Colors.Blue + "Unable to " + Colors.White + "%v cookies with error: " + Colors.Red + "%v\n" + Colors.Reset,
		CreateFile:   Colors.Blue + "Unable to create file " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Divider:      Colors.Grey + "*-------------------------------------------------------*\n" + Colors.Reset,
		Edit:         Colors.Blue + "Unable to edit " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Execute:      Colors.Blue + "Unable to execute template with error: " + Colors.Red + "%v\n" + Colors.Reset,
		KeyValuePair: Colors.Blue + "%v: " + Colors.White + "%v\n",
		Login:        Colors.Blue + "Unable to login " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		NoRows:       Colors.Blue + "No rows returned for " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v\n" + Colors.Reset,
		NotFound:     Colors.Blue + "Unable to find %v: " + Colors.White + "%v " + Colors.Blue + "called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Open:         Colors.Blue + "Unable to open " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error " + Colors.Red + "%v\n" + Colors.Reset,
		Parse:        Colors.Blue + "Unable to parse " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Post:         Colors.Blue + "Unable to post with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Printf:       Colors.Blue + "Unable to print with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Protected:    Colors.Blue + "CSRF validation failed for user " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Query:        Colors.Blue + "Unable to query " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Read:         Colors.Blue + "Unable to read " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		RetrieveFile: Colors.Blue + "Unable to retrieve file " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Register:     Colors.Blue + "Unable to register with error: " + Colors.Red + "%v\n" + Colors.Reset,
		SaveFile:     Colors.Blue + "Unable to save file " + Colors.White + "%v" + Colors.Blue + " to " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		Unmarshal:    Colors.Blue + "Unable to unmarshall " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v" + Colors.Blue + " with error: " + Colors.Red + "%v\n" + Colors.Reset,
		UserModel:    Colors.Blue + "Usermodel or DB called in " + Colors.White + "%v" + Colors.Blue + " for " + Colors.White + "%v" + Colors.Blue + " is nil\n" + Colors.Reset,
		Write:        Colors.Blue + "Unable to write to " + Colors.White + "%v" + Colors.Blue + " called by " + Colors.White + "%v\n" + Colors.Reset,
	}
	return customErrors
}
