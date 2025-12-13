// Run this demo with: go run cmd/demo_structured_logging.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gary-norman/forum/internal/models"
	"github.com/gary-norman/forum/internal/sqlite"
	"github.com/gary-norman/forum/internal/workers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to database
	db, err := sql.Open("sqlite3", "/var/lib/db-codex/dev_forum_database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create logger pool
	loggerPool := workers.NewLoggerPool(3, 100, db)
	loggerPool.Start()
	defer loggerPool.Shutdown(nil)

	// Simulate users
	user1 := models.NewUUIDField()
	user2 := models.NewUUIDField()
	user3 := models.NewUUIDField()

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  DEMO: Structured Logging with Searchable Context")
	fmt.Println("═══════════════════════════════════════════════════════════\n")

	// Scenario 1: User creates a post
	fmt.Println("📝 User 1 creates a post about Go concurrency...")
	loggerPool.LogPostCreated(
		user1,
		101, // postID
		5,   // channelID
		"Advanced Go Concurrency Patterns",
		[]string{"golang", "concurrency", "channels", "worker-pool"},
	)

	// Scenario 2: Multiple users view the post
	fmt.Println("👁️  Users viewing the post...")
	loggerPool.LogPostViewed(user1, 101, 5)
	loggerPool.LogPostViewed(user2, 101, 5)
	loggerPool.LogPostViewed(user3, 101, 5)

	// Scenario 3: User comments on the post
	fmt.Println("💬 User 2 comments on the post...")
	loggerPool.LogCommentCreated(user2, 201, 101, 150) // 150 char comment

	// Scenario 4: User reacts to the post
	fmt.Println("👍 User 3 likes the post...")
	loggerPool.LogReaction(user3, 101, "like")

	// Scenario 5: User joins a channel
	fmt.Println("📢 User 2 joins golang-discussions channel...")
	loggerPool.LogChannelJoined(user2, 5, "golang-discussions")

	// Scenario 6: User searches
	fmt.Println("🔍 User 3 searches for 'worker pool'...")
	loggerPool.LogSearchQuery(user3, "worker pool", 8) // 8 results found

	// Scenario 7: User login
	fmt.Println("🔐 User 1 logs in...")
	loggerPool.LogUserLogin(user1, "alice", "192.168.1.100", "Mozilla/5.0")

	// Wait for async writes
	time.Sleep(2 * time.Second)

	fmt.Println("\n" + "═══════════════════════════════════════════════════════════")
	fmt.Println("  SEARCHING STRUCTURED LOGS")
	fmt.Println("═══════════════════════════════════════════════════════════\n")

	loggingModel := &sqlite.LoggingModel{DB: db}

	// Query 1: Get all activity by user 1
	fmt.Println("🔎 Query 1: All activity by User 1")
	fmt.Println("───────────────────────────────────────────────────────────")
	activities, err := loggingModel.GetUserActivity(user1.String(), 50)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for _, activity := range activities {
			fmt.Printf("  [%s] %s - %s (%s #%d)\n",
				activity.Timestamp.Format("15:04:05"),
				activity.Action,
				activity.Description,
				activity.EntityType,
				activity.EntityID,
			)
		}
	}

	// Query 2: Get interactions with post 101
	fmt.Println("\n🔎 Query 2: All interactions with Post #101")
	fmt.Println("───────────────────────────────────────────────────────────")
	interactions, err := loggingModel.GetPostInteractions(101)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for _, interaction := range interactions {
			fmt.Printf("  [%s] User %s - %s\n",
				interaction.Timestamp.Format("15:04:05"),
				interaction.UserID[:8]+"...",
				interaction.Description,
			)
		}
	}

	// Query 3: Search for "golang" keyword
	fmt.Println("\n🔎 Query 3: Search for keyword 'golang'")
	fmt.Println("───────────────────────────────────────────────────────────")
	results, err := loggingModel.SearchByKeyword("golang", 10)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for _, result := range results {
			fmt.Printf("  [%s] %s - User %s\n",
				result.Timestamp.Format("15:04:05"),
				result.MetricName,
				result.UserID[:8]+"...",
			)
		}
	}

	// Query 4: Channel activity
	fmt.Println("\n🔎 Query 4: Activity in Channel #5")
	fmt.Println("───────────────────────────────────────────────────────────")
	channelActivity, err := loggingModel.GetChannelActivity(5, 50)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for _, activity := range channelActivity {
			fmt.Printf("  [%s] %s - %s by User %s\n",
				activity.Timestamp.Format("15:04:05"),
				activity.Action,
				activity.Description,
				activity.UserID[:8]+"...",
			)
		}
	}

	// Query 5: Most active users (last 7 days)
	fmt.Println("\n🔎 Query 5: Most Active Users (last 7 days)")
	fmt.Println("───────────────────────────────────────────────────────────")
	activeUsers, err := loggingModel.GetMostActiveUsers(7, 5)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for i, user := range activeUsers {
			fmt.Printf("  #%d: User %s... - %d actions (last: %s)\n",
				i+1,
				user.UserID[:8],
				user.ActivityCount,
				user.LastActivity.Format("15:04:05"),
			)
		}
	}

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  DEMO COMPLETE")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("\n💡 You can now query SystemMetrics table with JSON functions:")
	fmt.Println("   - Filter by user_id, post_id, channel_id")
	fmt.Println("   - Search keywords array")
	fmt.Println("   - Query metadata for custom fields")
	fmt.Println("   - Build activity feeds, analytics dashboards, etc.\n")
}
