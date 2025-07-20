package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gary-norman/forum/internal/app"
	mw "github.com/gary-norman/forum/internal/http/middleware"
	"github.com/gary-norman/forum/internal/models"
	"github.com/gary-norman/forum/internal/view"
)

type PostHandler struct {
	App      *app.App
	Channel  *ChannelHandler
	Comment  *CommentHandler
	Reaction *ReactionHandler
}

// GetUserPosts returns a slice of Posts that belong to channels the user follows. If no user is logged in, it returns all posts
func (p *PostHandler) GetUserPosts(user *models.User, allPosts []models.Post) []models.Post {
	if user == nil {
		return allPosts
	}
	var ownedAndJoinedChannels []models.Channel
	var postsInUserChannels []models.Post

	memberships, memberErr := p.App.Memberships.UserMemberships(user.ID)
	if memberErr != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "getHome > UserMemberships", memberErr)
	}
	ownedChannels, ownedErr := p.App.Channels.OwnedOrJoinedByCurrentUser(user.ID)
	if ownedErr != nil {
		log.Printf(ErrorMsgs().Query, "GetUserPosts > user owned channels", ownedErr)
	}
	joinedChannels, joinedErr := p.Channel.JoinedByCurrentUser(memberships)
	if joinedErr != nil {
		log.Printf(ErrorMsgs().Query, "user joined channels", joinedErr)
	}
	ownedAndJoinedChannels = append(ownedChannels, joinedChannels...)

	if len(ownedAndJoinedChannels) == 0 {
		return allPosts
	}

	// Create a set of user's channel IDs for efficient lookup
	userChannelIDSet := make(map[int64]bool)
	for _, membership := range ownedAndJoinedChannels {
		userChannelIDSet[membership.ID] = true
	}

	// Iterate over all the posts and check if their ChannelID is in the user's channel set
	for _, post := range allPosts {
		if _, exists := userChannelIDSet[post.ChannelID]; exists {
			postsInUserChannels = append(postsInUserChannels, post)
		}
	}

	return postsInUserChannels
}

func (p *PostHandler) GetThisPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var thisPost models.Post
	var posts []models.Post
	isMember := false
	var isMemberErr error
	isOwner := false

	userLoggedIn := true
	currentUser, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		fmt.Printf(ErrorMsgs().NotFound, "currentUser", "getThisPost", "_")
		userLoggedIn = false
	}

	// Parse post ID from URL
	postID, err := models.GetIntFromPathValue(r.PathValue("postId"))
	if err != nil {
		http.Error(w, `{"error": "invalid post ID"}`, http.StatusBadRequest)
	}

	// Fetch the post
	post, err := p.App.Posts.GetPostByID(postID)
	if err != nil {
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
	}
	posts = append(posts, post)
	foundPosts, err := p.Comment.GetPostsComments(posts)
	if err != nil {
		http.Error(w, `{"error": "error fetching post comments"}`, http.StatusInternalServerError)
	}
	foundPosts = p.Reaction.GetPostsLikesAndDislikes(foundPosts)
	thisPost = foundPosts[0]

	thisPost.ChannelID, thisPost.ChannelName, err = p.Channel.GetChannelInfoFromPostID(thisPost.ID)
	if err != nil {
		http.Error(w, `{"error": "error fetching channel info"}`, http.StatusInternalServerError)
	}

	models.UpdateTimeSince(&thisPost)

	// Fetch the channel
	channel, err := p.App.Channels.GetChannelByID(thisPost.ChannelID)
	if err != nil {
		log.Printf(ErrorMsgs().Query, "GetThisPost > GetChannelByID", err)
	}

	// Fetch the author
	author, err := p.App.Users.GetUserByUsername(thisPost.Author, "GetThisPost")
	if err != nil {
		log.Printf(ErrorMsgs().Query, "GetThisPost > GetUserByUsername", err)
	}

	if userLoggedIn {
		currentUser.Followers, currentUser.Following, err = p.App.Loyalty.CountUsers(currentUser.ID)
		if err != nil {
			http.Error(w, `{"error": "error fetching user loyalty"}`, http.StatusInternalServerError)
		}
		// Fetch if the user is a member of the channel
		isMember, isMemberErr = p.App.Channels.IsUserMemberOfChannel(currentUser.ID, channel.ID)
		if isMemberErr != nil {
			log.Printf(ErrorMsgs().Query, "GetThisPost > IsUserMemberOfChannel", err)
		}
		// Fetch if the user is the owner of the channel
		isOwner = currentUser.ID == channel.OwnerID
	}

	data := models.PostPage{
		UserID:      models.NewUUIDField(), // Default value of 0 for logged out users
		CurrentUser: currentUser,
		Instance:    "post-page",
		ThisPost:    thisPost,
		Author:      author,
		ThisChannel: channel,
		OwnerName:   author.Username,
		IsOwned:     isOwner,
		IsJoined:    isMember,
		ImagePaths:  p.App.Paths,
	}
	view.RenderPageData(w, data)
}

func (p *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	tpl, parseErr := template.ParseFiles("./assets/templates/posts.create.html")
	if parseErr != nil {
		http.Error(w, parseErr.Error(), 500)
		log.Printf(ErrorMsgs().Parse, "./assets/templates/posts.create.html", "createPost", parseErr)
		return
	}

	execErr := tpl.Execute(w, nil)
	if execErr != nil {
		log.Printf(ErrorMsgs().Execute, execErr)
		return
	}
}

func (p *PostHandler) StorePost(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	parseErr := r.ParseMultipartForm(10 << 20)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), 400)
		log.Printf(ErrorMsgs().Parse, "storePost", parseErr)
		return
	}
	channels := r.Form["channel[]"] // Retrieve the slice of checked channel IDs
	// get channel name

	// SECTION getting channel data (for reverting to single channel post
	//selectionJSON := r.PostForm.Get("channel")
	//if selectionJSON == "" {
	//	http.Error(w, "No selection provided", http.StatusBadRequest)
	//	return
	//}
	//var channelData models.ChannelData
	//if err := json.Unmarshal([]byte(selectionJSON), &channelData); err != nil {
	//	log.Printf(ErrorMsgs().Unmarshal, selectionJSON, err)
	//	http.Error(w, "Invalid selection format", http.StatusBadRequest)
	//	return
	//}
	//fmt.Printf(ErrorMsgs().KeyValuePair, "channelName", channelData.ChannelName)
	//fmt.Printf(ErrorMsgs().KeyValuePair, "channelID", channelData.ChannelID)
	//fmt.Printf(ErrorMsgs().KeyValuePair, "commentable", r.PostForm.Get("commentable"))

	createPostData := models.Post{
		Title:         r.PostForm.Get("title"),
		Content:       r.PostForm.Get("content"),
		Images:        "",
		Author:        user.Username,
		AuthorID:      user.ID,
		AuthorAvatar:  user.Avatar,
		IsCommentable: false,
		IsFlagged:     false,
	}

	fmt.Printf(ErrorMsgs().KeyValuePair, "Title", createPostData.Title)
	fmt.Printf(ErrorMsgs().KeyValuePair, "Content", createPostData.Content)
	fmt.Printf(ErrorMsgs().KeyValuePair, "Images", createPostData.Images)
	fmt.Printf(ErrorMsgs().KeyValuePair, "Author", createPostData.Author)
	fmt.Printf(ErrorMsgs().KeyValuePair, "Username", user.Username)

	fmt.Printf(ErrorMsgs().KeyValuePair, "AuthorID", createPostData.AuthorID)
	fmt.Printf(ErrorMsgs().KeyValuePair, "AuthorID type", reflect.TypeOf(createPostData.AuthorID))
	fmt.Printf(ErrorMsgs().KeyValuePair, "ID", user.ID)
	fmt.Printf(ErrorMsgs().KeyValuePair, "ID type", reflect.TypeOf(user.ID))

	fmt.Printf(ErrorMsgs().KeyValuePair, "AuthorAvatar", createPostData.AuthorAvatar)
	fmt.Printf(ErrorMsgs().KeyValuePair, "Avatar", user.Avatar)

	fmt.Printf(ErrorMsgs().KeyValuePair, "IsComm", createPostData.IsCommentable)
	fmt.Printf(ErrorMsgs().KeyValuePair, "IsFlag", createPostData.IsFlagged)

	if r.PostForm.Get("commentable") == "on" {
		createPostData.IsCommentable = true
	}
	createPostData.Images = GetFileName(r, "file-drop", "storePost", "post")
	/*createPostData.ChannelName = channelData.ChannelName
	createPostData.ChannelID, _ = strconv.Atoi(channelData.ChannelID)*/

	postID, insertErr := p.App.Posts.Insert(
		createPostData.Title,
		createPostData.Content,
		createPostData.Images,
		createPostData.Author,
		createPostData.AuthorAvatar,
		createPostData.AuthorID,
		createPostData.IsCommentable,
		createPostData.IsFlagged,
	)

	if insertErr != nil {
		log.Printf(ErrorMsgs().Post, insertErr)
		http.Error(w, insertErr.Error(), 500)
		return
	}

	for c := range channels {
		channelID, convErr := strconv.ParseInt(channels[c], 10, 64)
		if convErr != nil {
			log.Printf(ErrorMsgs().Convert, channels[c], "StorePost > GetChannelID", convErr)
			log.Printf("Unable to convert %v to integer\n", channels[c])
			continue
		}
		postToChannelErr := p.App.Channels.AddPostToChannel(channelID, postID)
		if postToChannelErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "channelID", channelID)
			log.Printf(ErrorMsgs().KeyValuePair, "postID", postID)
			log.Printf(ErrorMsgs().KeyValuePair, "postToChannelErr", postToChannelErr)
			http.Error(w, postToChannelErr.Error(), 500)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
