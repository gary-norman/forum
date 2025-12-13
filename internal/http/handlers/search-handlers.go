package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gary-norman/forum/internal/app"
	mw "github.com/gary-norman/forum/internal/http/middleware"
)

type SearchHandler struct {
	App *app.App
}

func (s *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Use concurrent search with request context
	result, err := ConcurrentSearch(r.Context(), s.App)
	if err != nil {
		log.Printf("Search completed with errors: %v", err)
		// Continue even with partial errors - result may still have data
	}

	// Enrich posts with channel information
	enrichedPosts := enrichPostsWithChannels(s.App, result.Posts, result.Channels)

	currentUser, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		log.Printf(ErrorMsgs.KeyValuePair, "User is not logged in. CurrentUser", currentUser)
	}

	searchResults := map[string]any{
		"users":    result.Users,
		"channels": result.Channels,
		"posts":    enrichedPosts,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(searchResults); err != nil {
		log.Printf(ErrorMsgs.Encode, "search results", err)
		http.Error(w, "Error encoding search results", http.StatusInternalServerError)
		return
	}

	log.Printf("[GET] /search - 200 (%dms)", time.Since(start).Milliseconds())
}
