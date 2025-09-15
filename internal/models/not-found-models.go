package models

type NotFound struct {
	Instance string
	Location string
	Message  string
}

func (n NotFound) GetInstance() string { return n.Instance }
