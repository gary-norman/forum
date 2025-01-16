package models

type Errors struct {
	Open      string
	Read      string
	Write     string
	Close     string
	ConnInit  string
	ConnConn  string
	ConnClose string
	Parse     string
	Execute   string
	Post      string
	Delete    string
	Update    string
	Insert    string
}

func CreateErrorMessages() *Errors {
	errors := &Errors{
		Open:      "Unable to open %v called by %v\n",
		Read:      "Unable to read %v called by %v\n",
		Write:     "Unable to write to %v called by %v\n",
		Close:     "Unable to close %v called by %v\n",
		ConnInit:  "Unable to initialise connection %v called by %v\n",
		ConnConn:  "Unable to connect to %v called by %v\n",
		ConnClose: "Unable to close connection to %v called by %v\n",
		Parse:     "Unable to parse %v called by %v with error: %v\n",
		Execute:   "Unable to execute template with error: %v\n",
		Post:      "Unable to post with error: %v\n",
		Delete:    "Unable to delete with error: %v\n",
		Update:    "Unable to update with error: %v\n",
		Insert:    "Unable to insert with error: %v\n",
	}
	return errors
}
