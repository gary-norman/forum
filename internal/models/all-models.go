package models

type All struct {
	Users       []User
	Posts       []Post
	Channels    []Channel
	Comments    []Comment
	Images      []Image
	Reactions   []Reaction
	Loyalty     []Loyalty
	Rules       []Rule
	Memberships []Membership
	Mods        []Mod
}
