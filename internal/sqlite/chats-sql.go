package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/gary-norman/forum/internal/models"
)

type ChatModel struct {
	DB *sql.DB
}

func (c *ChatModel) CreateChat(chatType, name string, groupID, buddyID models.UUIDField) (int64, error) {
	query := "INSERT INTO Chats (Type, Name, GroupID, BuddyID, Created) VALUES (?, ?, ?, ?, DateTime('now'))"
	result, err := c.DB.Exec(query, chatType, name, groupID, buddyID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert chat: %w", err)
	}
	chatID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get chat ID: %w", err)
	}

	return int64(chatID), nil
}

func (c *ChatModel) CreateMessage(userID models.UUIDField, message string) (int64, error) {
	query := "INSERT INTO Messages (ChatID, UserID, Created, Content) VALUES (?, ?, DateTime('now'), ?)"
	result, err := c.DB.Exec(query, userID, message)
	if err != nil {
		return 0, fmt.Errorf("failed to insert message: %w", err)
	}
	messageID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get message ID: %w", err)
	}

	return int64(messageID), nil
}

func (c *ChatModel) InsertMessageIntoChats(chatID, messageID models.UUIDField) error {
	query := "INSERT INTO ChatMessages (ChatID, MessageID) VALUES (?, ?)"
	_, err := c.DB.Exec(query, chatID, messageID)
	if err != nil {
		return fmt.Errorf("failed to insert message into chat: %w", err)
	}

	return nil
}

func (c *ChatModel) AttachUserToChat(chatID, userID models.UUIDField) error {
	query := "INSERT INTO ChatUsers (ChatID, UserID) VALUES (?, ?)"
	_, err := c.DB.Exec(query, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to attach user to chat: %w", err)
	}

	return nil
}

func (c *ChatModel) GetUserChats(userID models.UUIDField) ([]models.UUIDField, error) {
	query := `SELECT ID FROM ChatUsers WHERE UserID = ?`
	rows, err := c.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user chats: %w", err)
	}
	defer rows.Close()

	var chatIDs []models.UUIDField
	for rows.Next() {
		var chatID models.UUIDField
		if err := rows.Scan(&chatID); err != nil {
			return nil, fmt.Errorf("failed to scan chat ID: %w", err)
		}
		chatIDs = append(chatIDs, chatID)
	}

	return chatIDs, nil
}

func (c *ChatModel) GetChats(chatID models.UUIDField) ([]models.Chat, error) {
	query := "SELECT ID, Type, Created, LastActive, GroupID, BuddyID FROM Chats WHERE ID = ?"
	rows, err := c.DB.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chats: %w", err)
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.ChatType, &chat.Created, &chat.LastActive, &chat.Group.ID, &chat.Buddy.ID); err != nil {
			return nil, fmt.Errorf("failed to scan chat: %w", err)
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
