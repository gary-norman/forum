package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gary-norman/forum/internal/app"
	"github.com/gary-norman/forum/internal/models"
)

type ReactionHandler struct {
	App *app.App
}

// getPostsLikesAndDislikes updates the reactions of each post in the given slice
func (h *ReactionHandler) GetPostsLikesAndDislikes(posts []models.Post) []models.Post {
	for p, post := range posts {
		likes, dislikes, err := h.App.Reactions.CountReactions(post.ID, 0) // Pass 0 for CommentID if it's a post
		// fmt.Printf("PostID: %v, Likes: %v, Dislikes: %v\n", posts[i].ID, likes, dislikes)
		if err != nil {
			log.Printf("Error counting reactions: %v", err)
			likes, dislikes = 0, 0 // Default values if there is an error
		}
		models.React(&posts[p], likes, dislikes)
	}
	return posts
}

// getCommentsLikesAndDislikes updates the reactions of each comment in the given slice
func (h *ReactionHandler) GetCommentsLikesAndDislikes(comments []models.Comment) []models.Comment {
	for c, comment := range comments {
		likes, dislikes, likesErr := h.App.Reactions.CountReactions(0, comment.ID) // Pass 0 for PostID if it's a comment
		// fmt.Printf("PostID: %v, Likes: %v, Dislikes: %v\n", posts[i].ID, likes, dislikes)
		if likesErr != nil {
			log.Printf("Error counting reactions: %v", likesErr)
			likes, dislikes = 0, 0 // Default values if there is an error
		}
		models.React(&comments[c], likes, dislikes)
	}
	return comments
}

func (h *ReactionHandler) StoreReaction(w http.ResponseWriter, r *http.Request) {
	// log.Printf("using storeReaction()")

	// Variable to hold the decoded data
	var reactionData models.Reaction

	// Decode the JSON request body into the reactionData variable
	err := json.NewDecoder(r.Body).Decode(&reactionData)
	fmt.Printf("reactionData received: %v\n", &reactionData)
	if err != nil {
		// Use the error message from the Errors struct for decoding errors
		log.Println("Error decoding JSON")
		log.Printf(ErrorMsgs().Parse, "storeReaction", err)
		// TODO expect: 100-continue on the request header (for all of these)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//// Validate that at least one of reactedPostID or reactedCommentID is non-zero
	if (reactionData.ReactedPostID == nil || *reactionData.ReactedPostID == 0) && (reactionData.ReactedCommentID == nil || *reactionData.ReactedCommentID == 0) {
		log.Println("You must react to either a post or a comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postID, commentID := 0, 0

	if reactionData.ReactedPostID != nil {
		// log.Println("ReactedPostID:", *reactionData.ReactedPostID)
		postID = *reactionData.ReactedPostID
	}

	if reactionData.ReactedCommentID != nil {
		// log.Printf("ReactedCommentID: %d", *reactionData.ReactedPostID)
		commentID = *reactionData.ReactedCommentID
	}

	fmt.Printf("commentID after conversion: %v\n", commentID)
	fmt.Printf("postID after conversion: %v\n", postID)

	// Check if the user already reacted (like/dislike) and update or delete the reaction if needed
	existingReaction, err := h.App.Reactions.CheckExistingReaction(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
	if err != nil {
		// Use your custom error message for fetching errors
		log.Printf(ErrorMsgs().Read, "storeReaction > h.App.reactions.CheckExistingReaction", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// If there is an existing reaction, toggle it (i.e., remove it if the user reacts again to the same thing)
	if existingReaction != nil {
		log.Printf("Existing Reaction: %+v", existingReaction)

		// If the user likes a post or comment again, remove the like/dislike (toggle behavior)
		if existingReaction.Liked == reactionData.Liked && existingReaction.Disliked == reactionData.Disliked {
			// Reaction is the same, so remove it
			err = h.App.Reactions.Delete(existingReaction.ID)
			if err != nil {
				// Use your custom error message for deletion errors
				log.Printf(ErrorMsgs().Delete, "storeReaction > h.App.reactions.Delete", err)
				w.WriteHeader(http.StatusInternalServerError)
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
		err = h.App.Reactions.Update(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
		if err != nil {
			// Use your custom error message for update errors
			log.Printf(ErrorMsgs().Delete, "storeReaction > h.App.reactions.Update", err)
			w.WriteHeader(http.StatusInternalServerError)
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
	err = h.App.Reactions.Upsert(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
	fmt.Printf("Upserting liked: %v, disliked: %v, authorID: %v, postID: %v, commentID: %v\n", reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
	if err != nil {
		log.Printf(ErrorMsgs().Delete, "storeReaction > h.App.reactions.Upsert", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Check reaction and store it in the database, or handle errors
	// Respond with a JSON response
	w.Header().Set("Content-Type", "application/json")
	// Send a response indicating success
	// w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction added (in go)"})
	if err != nil {
		log.Printf(ErrorMsgs().Post, err)
		http.Error(w, err.Error(), 500)
		return
	}
	// http.Redirect(w, r, "/", http.StatusFound)
}
