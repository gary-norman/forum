package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gary-norman/forum/internal/models"
)

func (app *app) getThisPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var thisPost models.Post
	var posts []models.Post

	// userLoggedIn := true
	// currentUser, ok := getUserFromContext(r.Context())
	// if !ok {
	// 	fmt.Printf(ErrorMsgs().NotFound, "currentUser", "getThisPost", "_")
	// 	userLoggedIn = false
	// }

	// Parse post ID from URL
	postId, err := models.GetIntFromPathValue(r.PathValue("postId"))
	if err != nil {
		http.Error(w, `{"error": "invalid post ID"}`, http.StatusBadRequest)
	}

	// Fetch the post
	post, err := app.posts.GetPostByID(postId)
	if err != nil {
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
	}
	posts = append(posts, post)
	foundPosts, err := app.getPostsComments(posts)
	if err != nil {
		http.Error(w, `{"error": "error fetching post comments"}`, http.StatusInternalServerError)
	}
	foundPosts = app.getPostsLikesAndDislikes(foundPosts)
	thisPost = foundPosts[0]

	thisPost.ChannelID, thisPost.ChannelName, err = app.GetChannelInfoFromPostID(thisPost.ID)
	if err != nil {
		http.Error(w, `{"error": "error fetching channel info"}`, http.StatusInternalServerError)
	}

	models.UpdateTimeSince(&thisPost)

	// if userLoggedIn {
	// 	currentUser.Followers, currentUser.Following, err = app.loyalty.CountUsers(currentUser.ID)
	// 	if err != nil {
	//      http.Error(w, `{"error": "error fetching user loyalty"}`, http.StatusInternalServerError)
	// 	}
	// }

	TemplateData.ThisPost = thisPost
	// TemplateData.CurrentUser = currentUser

	fmt.Printf(ErrorMsgs().KeyValuePair, "TemplateData.ThisPost", TemplateData.ThisPost.Title)

	response := map[string]any{
		"post":        thisPost.Title,
		"content":     thisPost.Content,
		"likes":       thisPost.Likes,
		"dislikes":    thisPost.Dislikes,
		"channelID":   thisPost.ChannelID,
		"channelName": thisPost.ChannelName,
		"timeSince":   thisPost.TimeSince,
	}

	// Write the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "error encoding JSON"}`, http.StatusInternalServerError)
	}
}

func (app *app) createPost(w http.ResponseWriter, r *http.Request) {
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

func (app *app) storePost(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
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
	fmt.Printf(ErrorMsgs().KeyValuePair, "commentable", r.PostForm.Get("commentable"))

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
	fmt.Printf(ErrorMsgs().KeyValuePair, "authorAvatar", createPostData.AuthorAvatar)
	if r.PostForm.Get("commentable") == "on" {
		createPostData.IsCommentable = true
	}
	createPostData.Images = GetFileName(r, "file-drop", "storePost", "post")
	/*createPostData.ChannelName = channelData.ChannelName
	createPostData.ChannelID, _ = strconv.Atoi(channelData.ChannelID)*/

	postID, insertErr := app.posts.Insert(
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
		channelID, convErr := strconv.Atoi(channels[c])
		if convErr != nil {
			log.Printf(ErrorMsgs().Convert, channels[c], "StorePost > GetChannelID", convErr)
			log.Printf("Unable to convert %v to integer\n", channels[c])
			continue
		}
		postToChannelErr := app.channels.AddPostToChannel(channelID, postID)
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
