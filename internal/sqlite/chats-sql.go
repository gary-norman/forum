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

func (c *ChatModel) CreateChatMessage(userID models.UUIDField, message string) (int64, error) {
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

func (c *ChatModel) GetUserChatIDs(userID models.UUIDField) ([]models.UUIDField, error) {
	query := `SELECT ID FROM ChatUsers WHERE UserID = ?`
	rows, err := c.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user chat IDs: %w", err)
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

func (c *ChatModel) GetChatMessages(chatID models.UUIDField) ([]models.ChatMessage, error) {
	query := `
		SELECT
			m.ID, m.ChatID, m.Created, m.Content,
			u.ID, u.Username, u.EmailAddress, u.Avatar, u.Banner,
			u.Description, u.Usertype, u.Created, u.Updated, u.IsFlagged,
			u.SessionToken, u.CSRFToken, u.HashedPassword
		FROM Messages m
		LEFT JOIN Users u ON m.UserID = u.ID
		WHERE m.ChatID = ?
		ORDER BY m.Created ASC
	`

	rows, err := c.DB.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat messages: %w", err)
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var message models.ChatMessage
		var user models.User

		// Use sql.Null types for potentially NULL user fields
		var (
			userID          sql.NullString
			username        sql.NullString
			email           sql.NullString
			avatar          sql.NullString
			banner          sql.NullString
			description     sql.NullString
			usertype        sql.NullString
			userCreated     sql.NullTime
			userUpdated     sql.NullTime
			isFlagged       sql.NullBool
			sessionToken    sql.NullString
			csrfToken       sql.NullString
			hashedPassword  sql.NullString
		)

		err := rows.Scan(
			&message.ID, &message.ChatID, &message.Created, &message.Content,
			&userID, &username, &email, &avatar, &banner,
			&description, &usertype, &userCreated, &userUpdated, &isFlagged,
			&sessionToken, &csrfToken, &hashedPassword,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat message: %w", err)
		}

		// Only populate Sender if user exists (LEFT JOIN might return NULLs)
		if userID.Valid {
			user.ID = models.UUIDField(userID.String)
			user.Username = username.String
			user.Email = email.String
			user.Avatar = avatar.String
			user.Banner = banner.String
			user.Description = description.String
			user.Usertype = usertype.String
			user.Created = userCreated.Time
			user.Updated = userUpdated.Time
			user.IsFlagged = isFlagged.Bool
			user.SessionToken = sessionToken.String
			user.CSRFToken = csrfToken.String
			user.HashedPassword = hashedPassword.String

			models.UpdateTimeSince(&user)
			message.Sender = &user
		} else {
			message.Sender = nil
		}

		messages = append(messages, message)
	}

	return messages, nil
}
