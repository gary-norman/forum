package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gary-norman/forum/internal/models"
)

func (app *app) getThisPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var thisPost models.Post
	var posts []models.Post

	userLoggedIn := true
	currentUser, ok := getUserFromContext(r.Context())
	if !ok {
		fmt.Printf(ErrorMsgs().NotFound, "currentUser", "getThisPost", "_")
		userLoggedIn = false
	}

	postId, err := strconv.Atoi(r.PathValue("postId"))
	if err != nil {
		fmt.Printf(ErrorMsgs().KeyValuePair, "convert postID", err)
	}
	post, err := app.posts.GetPostByID(postId)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "getHome > thisPost", err)
	}
	posts = append(posts, post)
	foundPosts, err := app.getPostsComments(posts)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "getHome > thisPost comments", err)
	}
	thisPost = foundPosts[0]
}
