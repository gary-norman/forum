package models

type ErrorPageData struct {
	Instance string
	Location string
	Message  string
}

func (e ErrorPageData) GetInstance() string { return e.Instance }
