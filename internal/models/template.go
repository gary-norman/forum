package models

type TemplateData struct {
	// ---------- users ----------
	AllUsers    []User
	RandomUser  User
	CurrentUser *User
	// ---------- posts ----------
	UserID    int64
	Posts     []Post
	UserPosts []Post
	// ---------- channels ----------
	Channels               []Channel
	AllChannels            []Channel
	ThisChannel            Channel
	ThisChannelOwnerName   string
	IsJoinedOrOwned        bool
	ThisChannelIsOwned     bool
	ThisChannelRules       []Rule
	ThisChannelPosts       []Post
	OwnedChannels          []Channel
	JoinedChannels         []Channel
	OwnedAndJoinedChannels []Channel
	// ---------- misc ----------
	Instance  string
	Images    []Image
	Reactions []Reaction
	ThisPost  Post
	ThisUser  User
	ImagePaths
}

func (p TemplateData) GetInstance() string { return p.Instance }
