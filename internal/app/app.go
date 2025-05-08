package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gary-norman/forum/internal/db"
	"github.com/gary-norman/forum/internal/models"
	"github.com/gary-norman/forum/internal/sqlite"
)

var (
	dbType   string = "SQLite"
	dbDriver string = "sqlite3"
	// dbPath    string = "db/forum_database.db"
	// development db for testing new setup
	dbPath     string = "db/dev_forum_database.db"
	schemaPath string = "schema.sql"
	imagePath  string = "db/userdata/images/"
)

func ErrorMsgs() *models.Errors {
	return models.CreateErrorMessages()
}

func Colors() *models.Colors {
	return models.CreateColors()
}

type App struct {
	Db          *sql.DB // Store DB reference for cleanup
	Users       *sqlite.UserModel
	Posts       *sqlite.PostModel
	Reactions   *sqlite.ReactionModel
	Saved       *sqlite.SavedModel
	Mods        *sqlite.ModModel
	Comments    *sqlite.CommentModel
	Images      *sqlite.ImageModel
	Channels    *sqlite.ChannelModel
	Flags       *sqlite.FlagModel
	Loyalty     *sqlite.LoyaltyModel
	Memberships *sqlite.MembershipModel
	Muted       *sqlite.MutedChannelModel
	Cookies     *sqlite.CookieModel
	Rules       *sqlite.RuleModel
	Paths       models.ImagePaths
}

func NewApp(db *sql.DB) *App {
	return &App{
		Db:          db,
		Users:       &sqlite.UserModel{DB: db},
		Posts:       &sqlite.PostModel{DB: db},
		Reactions:   &sqlite.ReactionModel{DB: db},
		Saved:       &sqlite.SavedModel{DB: db},
		Mods:        &sqlite.ModModel{DB: db},
		Comments:    &sqlite.CommentModel{DB: db},
		Images:      &sqlite.ImageModel{DB: db},
		Channels:    &sqlite.ChannelModel{DB: db},
		Flags:       &sqlite.FlagModel{DB: db},
		Loyalty:     &sqlite.LoyaltyModel{DB: db},
		Memberships: &sqlite.MembershipModel{DB: db},
		Muted:       &sqlite.MutedChannelModel{DB: db},
		Cookies:     &sqlite.CookieModel{DB: db}, // Not in sqlite, handled separately
		Rules:       &sqlite.RuleModel{DB: db},

		Paths: models.ImagePaths{
			Channel: imagePath + "channel-images/",
			Post:    imagePath + "post-images/",
			User:    imagePath + "user-images/",
		},
	}
}

func InitializeApp() (*App, func(), error) {
	// Open database connection
	// db, err := sql.Open(dbDriver, dbPath)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// os.Remove(dbPath)

	db, err := db.InitDB(dbPath, schemaPath)
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		db.Close() // Close DB if ping fails
		return nil, nil, err
	}

	// Get version
	var dbVersion string
	if err = db.QueryRow("select sqlite_version()").Scan(&dbVersion); err != nil {
		log.Printf(Colors().Red+"Error fetching SQLite version: %v\n"+Colors().Reset, err)
	}
	fmt.Printf(ErrorMsgs().DbSuccess, dbType, dbVersion)

	// App instance with DB reference
	appInstance := NewApp(db)

	// Cleanup function to close DB connection
	cleanup := func() {
		log.Println("Closing database connection...")
		if err := db.Close(); err != nil {
			log.Println(Colors().Red+"Error closing DB: %v"+Colors().Reset, err)
		} else {
			fmt.Println(Colors().Green + "Database closed successfully." + Colors().Reset)
		}
	}

	return appInstance, cleanup, nil
}
