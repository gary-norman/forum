// Package models contains data models used across the application.
package models

type ErrorPage struct {
	Data   ErrorPageData
	Status int
	Error  error
}

type TemplateData struct {
	// ---------- users ----------
	UserID      UUIDField
	AllUsers    []User
	RandomUser  User
	CurrentUser *User
	// ---------- posts ----------
	Posts     []Post
	UserPosts []Post
	// ---------- channels ----------
	// Channels               []Channel
	AllChannels []Channel
	ThisChannel Channel
	// ThisChannelOwnerName   string
	// IsJoinedOrOwned        bool
	// ThisChannelIsOwned     bool
	ThisChannelRules []Rule
	// ThisChannelPosts       []Post
	OwnedChannels          []Channel
	JoinedChannels         []Channel
	OwnedAndJoinedChannels []Channel
	// ---------- misc ----------
	Instance string
	// Images    []Image
	// Reactions []Reaction
	// ThisPost  Post
	ThisUser User
	ImagePaths
	// ErrorPage ErrorPage
	ErrorPage bool
}

func (p TemplateData) GetInstance() string { return p.Instance }
