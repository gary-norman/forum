package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gary-norman/forum/internal/colors"
	"github.com/gary-norman/forum/internal/models"
)

var (
	seedColors, _ = colors.UseFlavor("Mocha")
)

func runSeed(db *sql.DB) error {
	// Check if seed data already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE Username = ?", "TheCodexDonkey").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing seed data: %w", err)
	}
	if count > 0 {
		fmt.Printf("%s⊘ Seed data already exists, skipping...%s\n", seedColors.Yellow, seedColors.Reset)
		return nil
	}

	fmt.Printf("%s> Seeding database with initial data...%s\n", seedColors.CodexPink, seedColors.Reset)

	// create variables
	userID := models.NewUUIDField()
	now := time.Now().UTC().Format(time.RFC3339)

	// Copy default images
	err = models.CopyFile(
		"default-images/codex.png",
		"db/userdata/images/channel-images/codex.png",
	)
	if err != nil {
		return err
	}

	err = models.CopyFile(
		"default-images/donkey.png",
		"db/userdata/images/user-images/donkey.png",
	)
	if err != nil {
		return err
	}

	// Insert user
	_, err = db.Exec(`
	INSERT INTO Users (ID, Username, EmailAddress, Avatar, Banner, Description, Usertype, Created, IsFlagged, SessionToken, CsrfToken, HashedPassword)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, "TheCodexDonkey", "donkey@codex.com", "donkey.png", "",
		"I'm such a friendly donkey, and I'm here to show you around the wonderful world of Codex. I've already added you to my channel, where you'll find updates, information and any cool stuff I feel like sharing. Welcome to Codex!",
		0, now, 0, "", "", "$2a$14$qK2P4N12utI8c4dPS6AMaueafDJygKtdVHVLgNVq2wJM5MW5xjdVm")
	if err != nil {
		return fmt.Errorf("failed to insert %v seed data: %w", "user", err)
	}

	// Insert channel
	res, err := db.Exec(`
	INSERT INTO Channels (OwnerID, Name, Description, Created, Avatar, Banner, Privacy, IsFlagged, IsMuted)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, "WelcomeToCodex",
		"Welcome to Codex! This channel will guide you through the forum, and give you updates, insights, and generally keep you up to date with everything that's going on.",
		now, "codex.png", "", 0, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to insert %v seed data: %w", "channel", err)
	}
	channelID, _ := res.LastInsertId()

	// Insert post
	res, err = db.Exec(`
	INSERT INTO Posts (Title, Content, Images, Created, Author, AuthorAvatar, AuthorID, IsCommentable, IsFlagged)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"Welcome to Codex!",
		"We are so glad you could join us! Since you're here, you probably already know this, but just in case, this is a place for us all to share everything we know about coding - our hopes, fears, plans, anxieties... Everything that makes up this world in which we find ourselves. Hopefully, Codex will make it less daunting, and a world we enjoy living in. So, welcome! We're so glad you're here.",
		"noimage", now, "TheCodexDonkey", "donkey.png", userID, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to insert %v seed data: %w", "post", err)
	}
	postID, _ := res.LastInsertId()

	// Attach post to channel
	_, err = db.Exec(`INSERT INTO PostChannels (PostID, ChannelID, Created) VALUES (?, ?, ?)`, postID, channelID, now)
	if err != nil {
		return fmt.Errorf("failed to insert %v seed data: %w", "post", err)
	}

	fmt.Printf("%s✓ Successfully seeded database%s\n", seedColors.CodexGreen, seedColors.Reset)
	return nil
}
