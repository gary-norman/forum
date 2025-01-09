package models

type Errors struct {
	Close     string
	ConnClose string
	ConnInit  string
	ConnConn  string
	Cookies   string
	Execute   string
	Login     string
	NoRows    string
	Open      string
	Parse     string
	Post      string
	Protected string
	Query     string
	Read      string
	Register  string
	UserModel string
	Write     string
}

func CreateErrorMessages() *Errors {
	errors := &Errors{
		Close:     "Unable to close %v called by %v\n",
		ConnConn:  "Unable to connect to %v called by %v\n",
		ConnClose: "Unable to close connection to %v called by %v\n",
		ConnInit:  "Unable to initialise connection %v called by %v\n",
		Cookies:   "Unable to set cookies with error: %v\n",
		Execute:   "Unable to execute template with error: %v\n",
		Login:     "Unable to login with error: %v\n",
		NoRows:    "No rows returned for %v called by %v\n",
		Open:      "Unable to open %v called by %v\n",
		Parse:     "Unable to parse %v called by %v with error: %v\n",
		Post:      "Unable to post with error: %v\n",
		Protected: "CSRF validation failed for user %v with error: %v\n",
		Query:     "Unable to query %v with error: %v\n",
		Read:      "Unable to read %v called by %v\n",
		Register:  "Unable to register with error: %v\n",
		UserModel: "Usermodel or DB called in %v for %v is nil\n",
		Write:     "Unable to write to %v called by %v\n",
	}
	return errors
}
