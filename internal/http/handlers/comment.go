package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gary-norman/forum/internal/app"
	mw "github.com/gary-norman/forum/internal/http/middleware"
	"github.com/gary-norman/forum/internal/models"
)

type CommentHandler struct {
	App      *app.App
	Reaction *ReactionHandler
}

func (h *CommentHandler) StoreComment(w http.ResponseWriter, r *http.Request) {
	// SECTION getting user
	user, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "No user found in context", http.StatusUnauthorized)
		return
	}

	// Check if the method is POST, otherwise return Method Not Allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// SECTION retrieving comment form data
	// Variable to hold the decoded data
	var commentData models.Comment

	parseErr := r.ParseMultipartForm(10 << 20)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), 400)
		log.Printf(ErrorMsgs().Parse, "storeComment", parseErr)
		return
	}

	// SECTION setting channel data
	// Get channel data
	selectionJSON := r.PostForm.Get("channel")
	if selectionJSON == "" {
		http.Error(w, "No selection provided for channel", http.StatusBadRequest)
		return
	}

	var channelData models.ChannelData
	if err := json.Unmarshal([]byte(selectionJSON), &channelData); err != nil {
		log.Printf(ErrorMsgs().Unmarshal, selectionJSON, err)
		http.Error(w, "Invalid selection format", http.StatusBadRequest)
		return
	}

	postIDStr := r.PostForm.Get("postID")
	commentIDStr := r.PostForm.Get("commentID")

	// Convert postIDStr to an integer
	postID, postConvErr := strconv.ParseInt(postIDStr, 10, 64)
	if postConvErr != nil {
		log.Printf("Error converting postID: %v", postID)
	}

	// Convert commentIDStr to an integer
	commentID, commentConvErr := strconv.ParseInt(commentIDStr, 10, 64)
	if commentConvErr != nil {
		log.Printf("Error converting commentID: %v", commentID)
	}

	// Assign the returned values
	commentData = models.Comment{
		Content:            r.PostForm.Get("content"),
		CommentedPostID:    &postID,
		CommentedCommentID: &commentID,
		IsCommentable:      true,
		IsFlagged:          false,
		IsReply:            false,
		Author:             user.Username,
		AuthorID:           user.ID,
		AuthorAvatar:       user.Avatar,
		ChannelName:        channelData.ChannelName,
		ChannelID:          0,
	}

	commentData.ChannelID, _ = strconv.ParseInt(channelData.ChannelID, 10, 64)

	// Log the values
	fmt.Printf("commentData.CommentedPostID: %v\n", *commentData.CommentedPostID)
	fmt.Printf("commentData.CommentedCommentID: %v\n", *commentData.CommentedCommentID)

	// Insert the comment
	insertErr := h.App.Comments.Upsert(commentData)

	if insertErr != nil {
		log.Printf(ErrorMsgs().Comment, insertErr)
		http.Error(w, insertErr.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *CommentHandler) GetPostsComments(posts []models.Post) ([]models.Post, error) {
	for p, post := range posts {
		comments, err := h.App.Comments.GetCommentByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		comments = h.Reaction.GetCommentsLikesAndDislikes(comments)
		/// Filter comments that belong to the current post based on the postID and CommentedPostID
		var postComments []models.Comment
		var commentsCount int
		for c, comment := range comments {
			models.UpdateTimeSince(&comments[c])
			// For each comment, recursively assign its replies
			commentWithReplies := h.GetRepliesForComment(comment)
			postComments = append(postComments, commentWithReplies)
			commentsCount = len(postComments)
		}
		posts[p].Comments = postComments
		posts[p].CommentsCount = commentsCount
	}
	return posts, nil
}

// getRepliesForComment Recursively fetches replies for each comment
func (h *CommentHandler) GetRepliesForComment(comment models.Comment) models.Comment {
	// Find replies to the current comment
	var replies []models.Comment
	comments, _ := h.App.Comments.GetCommentByCommentID(comment.ID)
	comments = h.Reaction.GetCommentsLikesAndDislikes(comments)
	for r, possibleReply := range comments {
		models.UpdateTimeSince(&comments[r])
		replyWithReplies := h.GetRepliesForComment(possibleReply) // Recursively get replies for this reply
		replies = append(replies, replyWithReplies)
	}
	if len(replies) > 0 {
		comment.Replies = replies
	}
	// Return the comment with its replies
	return comment
}
