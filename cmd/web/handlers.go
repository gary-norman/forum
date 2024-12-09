package main

import (
	"encoding/json"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func (app *app) getHome(w http.ResponseWriter, r *http.Request) {
	ErrorMsgs := models.CreateErrorMessages()
	posts, err := app.posts.All()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t, err := template.ParseFiles("./assets/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf(ErrorMsgs.Parse, "./assets/templates/index.html", "getHome", err)
		return
	}

	// Retrieve likes and dislikes for each post
	for i, post := range posts {
		likes, dislikes, err := app.reactions.CountReactions(post.ChannelID, post.ID, 0) // Pass 0 for CommentID if it's a post
		if err != nil {
			log.Printf("Error counting reactions: %v", err)
			likes, dislikes = 0, 0 // Default values if there is an error
		}

		posts[i].Likes = likes
		posts[i].Dislikes = dislikes
	}

	err = t.Execute(w, map[string]any{"Posts": posts})
	if err != nil {
		log.Print(ErrorMsgs.Execute, err)
		return
	}
}

func (app *app) createPost(w http.ResponseWriter, r *http.Request) {
	ErrorMsgs := models.CreateErrorMessages()
	t, err := template.ParseFiles("./assets/templates/posts.create.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf(ErrorMsgs.Parse, "./assets/templates/posts.create.html", "createPost", err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf(ErrorMsgs.Execute, err)
		return
	}
}

func (app *app) storePost(w http.ResponseWriter, r *http.Request) {
	ErrorMsgs := models.CreateErrorMessages()
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Printf(ErrorMsgs.Parse, "./assets/templates/posts.create.html", "storePost", err)
		return
	}

	// Get the 'channel' value as a string
	channelStr := r.PostForm.Get("channel")
	// Convert the string to an integer
	channel, err := strconv.Atoi(channelStr)
	if err != nil {
		http.Error(w, "You must be a member of this channel to do that.", http.StatusBadRequest)
		return
	}

	// Get the 'author' value as a string
	authorStr := r.PostForm.Get("author")
	// Convert the string to an integer
	author, err := strconv.Atoi(authorStr)
	if err != nil {
		http.Error(w, "You must be logged in to do that.", http.StatusBadRequest)
		return
	}

	type FormData struct {
		commentable bool
		images      string
	}
	formData := FormData{
		commentable: false,
		images:      "noimage",
	}
	if r.PostForm.Get("commentable") != "" {
		formData.commentable = true
	}
	images := r.PostForm.Get("images")
	if images != "" {
		formData.images = images
	}

	err = app.posts.Insert(
		r.PostForm.Get("title"),
		r.PostForm.Get("content"),
		formData.images,
		channel,
		author,
		formData.commentable,
	)

	if err != nil {
		log.Printf(ErrorMsgs.Post, err)
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *app) storeReaction(w http.ResponseWriter, r *http.Request) {

	//log.Printf("using storeReaction()")
	// Using custom error messages
	ErrorMsgs := models.CreateErrorMessages()

	// Check if the method is POST, otherwise return Method Not Allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Variable to hold the decoded data
	var reactionData models.Reaction

	// Decode the JSON request body into the reactionData variable
	err := json.NewDecoder(r.Body).Decode(&reactionData)
	if err != nil {
		// Use the error message from the Errors struct for decoding errors
		http.Error(w, fmt.Sprintf(ErrorMsgs.Parse, "storeReaction", err), http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	// Validate that at least one of reactedPostID or reactedCommentID is non-zero
	if (reactionData.ReactedPostID == nil || *reactionData.ReactedPostID == 0) && (reactionData.ReactedCommentID == nil || *reactionData.ReactedCommentID == 0) {
		http.Error(w, "You must react to either a post or a comment", http.StatusBadRequest)
		return
	}

	postID, commentID := 0, 0

	if reactionData.ReactedPostID != nil {
		//log.Println("ReactedPostID:", *reactionData.ReactedPostID)
		postID = *reactionData.ReactedPostID
	} else {
		postID = 0
		//log.Println("ReactedPostID is nil")
	}

	if reactionData.ReactedCommentID != nil {
		//log.Printf("ReactedCommentID: %d", *reactionData.ReactedPostID)
		commentID = *reactionData.ReactedCommentID
	} else {
		commentID = 0
		//log.Println("ReactedCommentID is nil")
	}

	//// Check if the user already reacted (like/dislike) and update or delete the reaction if needed
	existingReaction, err := app.reactions.CheckExistingReaction(reactionData.AuthorID, reactionData.ChannelID, postID, commentID)
	if err != nil {
		// Use your custom error message for fetching errors
		http.Error(w, fmt.Sprintf(ErrorMsgs.Read, "storeReaction", err), http.StatusInternalServerError)
		log.Printf("Error fetching existing reaction: %v", err)
		return
	}

	if existingReaction != nil {
		log.Printf("Existing Reaction: %+v", existingReaction)
	}

	// If there is an existing reaction, toggle it (i.e., remove it if the user reacts again to the same thing)
	if existingReaction != nil {
		// If the user likes a post or comment again, remove the like/dislike (toggle behavior)
		if existingReaction.Liked == reactionData.Liked && existingReaction.Disliked == reactionData.Disliked {
			// Reaction is the same, so remove it
			err = app.reactions.Delete(existingReaction.ID)
			if err != nil {
				// Use your custom error message for deletion errors
				http.Error(w, fmt.Sprintf(ErrorMsgs.Delete, "storeReaction", err), http.StatusInternalServerError)
				log.Printf("Error deleting reaction: %v", err)
				return
			}
			// Send response back after successful deletion
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction removed"})
			if err != nil {
				return
			}
			return
		}

		// Otherwise, update the existing reaction
		err = app.reactions.Update(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, reactionData.ChannelID, postID, commentID)
		if err != nil {
			// Use your custom error message for update errors
			http.Error(w, fmt.Sprintf(ErrorMsgs.Update, "storeReaction", err), http.StatusInternalServerError)
			log.Printf("Error updating reaction: %v", err)
			return
		}
		// Send response back after successful update
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction updated"})
		if err != nil {
			return
		}
		return
	}

	// If no existing reaction, insert a new reaction
	err = app.reactions.Upsert(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, reactionData.ChannelID, postID, commentID)
	if err != nil {
		// Use your custom error message for insertion errors
		http.Error(w, fmt.Sprintf(ErrorMsgs.Insert, "storeReaction", err), http.StatusInternalServerError)
		log.Printf("Error inserting reaction: %v", err)
		return
	}

	// Check reaction and store it in the database, or handle errors
	// Respond with a JSON response
	w.Header().Set("Content-Type", "application/json")

	// Send a response indicating success
	//w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction added (in go)"})

	if err != nil {
		log.Printf(ErrorMsgs.Post, err)
		http.Error(w, err.Error(), 500)
		return
	}

	//http.Redirect(w, r, "/", http.StatusFound)
}

//func (app *app) storeReaction(w http.ResponseWriter, r *http.Request) {
//	log.Printf("using storeReaction()")
//
//	if r.Method != http.MethodPost {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	var reactionData models.Reaction
//	err := json.NewDecoder(r.Body).Decode(&reactionData)
//	if err != nil {
//		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
//		return
//	}
//
//	// Check reaction and store it in the database, or handle errors
//	// Respond with a JSON response
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction updated (in go)"})
//	if err != nil {
//		return
//	}
//}
