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
			postStr := fmt.Sprintf("PostID: %v", post.ID)
			log.Printf(ErrorMsgs().Generic, "Error counting reactions", postStr, "GetPostsLikesAndDislikes", err)
			likes, dislikes = 0, 0 // Default values if there is an error
		}
		models.React(&posts[p], likes, dislikes)
	}
	return posts
}

// getCommentsLikesAndDislikes updates the reactions of each comment in the given slice
func (h *ReactionHandler) GetCommentsLikesAndDislikes(comments []models.Comment) []models.Comment {
	for c, comment := range comments {
		likes, dislikes, err := h.App.Reactions.CountReactions(0, comment.ID) // Pass 0 for PostID if it's a comment
		// fmt.Printf("PostID: %v, Likes: %v, Dislikes: %v\n", posts[i].ID, likes, dislikes)
		if err != nil {
			postStr := fmt.Sprintf("CommentID: %v", comment.ID)
			log.Printf(ErrorMsgs().Generic, "Error counting reactions", postStr, "GetPostsLikesAndDislikes", err)
			likes, dislikes = 0, 0 // Default values if there is an error
		}
		models.React(&comments[c], likes, dislikes)
	}
	return comments
}

func (h *ReactionHandler) StoreReaction(w http.ResponseWriter, r *http.Request) {
	log.Printf("using storeReaction()")

	// Variable to hold the decoded data
	var input models.ReactionInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Convert AuthorID string to UUIDField
	authorID, err := models.UUIDFieldFromString(input.AuthorID)
	if err != nil {
		http.Error(w, "Invalid authorId", http.StatusBadRequest)
		return
	}

	reactionData := models.Reaction{
		Liked:            input.Liked,
		Disliked:         input.Disliked,
		AuthorID:         authorID,
		ReactedPostID:    input.ReactedPostID,
		ReactedCommentID: input.ReactedCommentID,
	}

	//// Validate that at least one of reactedPostID or reactedCommentID is non-zero
	if (reactionData.ReactedPostID == nil || *reactionData.ReactedPostID == 0) && (reactionData.ReactedCommentID == nil || *reactionData.ReactedCommentID == 0) {
		log.Println("You must react to either a post or a comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var updatedID int64
	var updatedStr string

	if reactionData.ReactedPostID != nil {
		reactionData.PostID = *reactionData.ReactedPostID
		// log.Println("ReactedPostID:", *reactionData.ReactedPostID)
		updatedID = *reactionData.ReactedPostID
		updatedStr = "post"
	} else {
		reactionData.CommentID = *reactionData.ReactedCommentID
		// log.Printf("ReactedCommentID: %d", *reactionData.ReactedPostID)
		updatedID = *reactionData.ReactedCommentID
		updatedStr = "comment"
	}

	updatedPair := fmt.Sprintf("%s: %v", updatedStr, updatedID)
	log.Printf(ErrorMsgs().KeyValuePair, "Updating like for", updatedPair)

	if err := h.App.Reactions.Upsert(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, reactionData.PostID, reactionData.CommentID); err != nil {
		log.Printf(ErrorMsgs().Update, updatedPair, "storeReaction > h.App.reactions.Upsert", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updatedReactions, err := h.App.Reactions.All()
	if err != nil {
		log.Printf(ErrorMsgs().Read, updatedReactions, "storeReaction > h.App.reactions.All", err)
		return
	}

	for _, reaction := range updatedReactions {
		if reaction.ReactedPostID != nil {
			reaction.PostID = *reaction.ReactedPostID
			reaction.CommentID = 0
		} else {
			reaction.PostID = 0
			reaction.CommentID = *reaction.ReactedCommentID
		}
		reaction.ReactedCommentID = nil
		reaction.ReactedPostID = nil
		models.JsonPost(reaction)
		fmt.Println(ErrorMsgs().Divider)
	}
}

// Check if the user already reacted (like/dislike) and update or delete the reaction if needed
// 	existingReaction, err := h.App.Reactions.CheckExistingReaction(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
// 	if err != nil {
// 		log.Printf(ErrorMsgs().Read, reactionData, "storeReaction > h.App.reactions.CheckExistingReaction", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
//
// 	// If there is an existing reaction, toggle it (i.e., remove it if the user reacts again to the same thing)
// 	if existingReaction != nil {
// 		log.Printf("Existing Reaction: %+v", existingReaction)
//
// 		// If the user likes a post or comment again, remove the like/dislike (toggle behavior)
// 		if existingReaction.Liked == reactionData.Liked && existingReaction.Disliked == reactionData.Disliked {
// 			// Reaction is the same, so remove it
// 			err = h.App.Reactions.Delete(existingReaction.ID)
// 			if err != nil {
// 				// Use your custom error message for deletion errors
// 				log.Printf(ErrorMsgs().Delete, updatedPair, "storeReaction > h.App.reactions.Delete", err)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}
// 			// Send response back after successful deletion
// 			w.WriteHeader(http.StatusOK)
// 			err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction removed"})
// 			if err != nil {
// 				return
// 			}
// 			return
// 		}
//
// 		// Otherwise, update the existing reaction
// 		err = h.App.Reactions.Update(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
// 		if err != nil {
// 			// Use your custom error message for update errors
// 			log.Printf(ErrorMsgs().Update, updatedPair, "storeReaction > h.App.reactions.Update", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		// Send response back after successful update
// 		w.WriteHeader(http.StatusOK)
// 		err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction updated"})
// 		if err != nil {
// 			return
// 		}
// 		return
// 	}
//
// 	// If no existing reaction, insert a new reaction
// 	err = h.App.Reactions.Upsert(reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
// 	fmt.Printf("Upserting liked: %v, disliked: %v, authorID: %v, postID: %v, commentID: %v\n", reactionData.Liked, reactionData.Disliked, reactionData.AuthorID, postID, commentID)
// 	if err != nil {
// 		log.Printf(ErrorMsgs().Update, updatedPair, "storeReaction > h.App.reactions.Upsert", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	// Check reaction and store it in the database, or handle errors
// 	// Respond with a JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	// Send a response indicating success
// 	// w.WriteHeader(http.StatusCreated)
// 	err = json.NewEncoder(w).Encode(map[string]string{"message": "Reaction added (in go)"})
// 	if err != nil {
// 		log.Printf(ErrorMsgs().Post, err)
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	// http.Redirect(w, r, "/", http.StatusFound)
// }
