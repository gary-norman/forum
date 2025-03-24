package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gary-norman/forum/internal/models"
)

type ChannelModel struct {
	DB *sql.DB
}

func (m *ChannelModel) Insert(ownerID int, name, description, avatar, banner string, privacy, isFlagged, isMuted bool) error {
	stmt := "INSERT INTO Channels (OwnerID, Name, Description, Created, Avatar, Banner, Privacy, IsFlagged, IsMuted) VALUES (?, ?, ?, DateTime('now'), ?, ?, ?, ?, ?)"
	_, err := m.DB.Exec(stmt, ownerID, name, description, avatar, banner, privacy, isFlagged, isMuted)
	return err
}

func (m *ChannelModel) OwnedOrJoinedByCurrentUser(ID int, column string) ([]models.Channel, error) {
	validColumns := map[string]bool{
		"ID":          true,
		"OwnerID":     true,
		"Name":        true,
		"Avatar":      true,
		"Banner":      true,
		"Description": true,
		"Created":     true,
		"Privacy":     true,
		"IsFlagged":   true,
		"IsMuted":     true,
	}

	if !validColumns[column] {
		return nil, fmt.Errorf("invalid column name: %s", column)
	}

	// Safely construct the query
	stmt := fmt.Sprintf(
		"SELECT ID, OwnerID, Name, Avatar, Banner, Description, Created, Privacy, IsFlagged, IsMuted FROM Channels WHERE %s = ?",
		column,
	)
	rows, queryErr := m.DB.Query(stmt, ID)
	if queryErr != nil {
		return nil, queryErr
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "OwnedOrJoinedByCurrentUser", closeErr)
		}
	}()

	var channels []models.Channel
	for rows.Next() {
		var c models.Channel
		scanErr := rows.Scan(&c.ID, &c.OwnerID, &c.Name, &c.Avatar, &c.Banner, &c.Description, &c.Created, &c.Privacy, &c.IsFlagged, &c.IsMuted)
		if scanErr != nil {
			return nil, scanErr
		}
		if column == "OwnerID" {
			fmt.Printf(ErrorMsgs().KeyValuePair, "updating Owned of", c.Name)
			c.Owned = true
			fmt.Printf(ErrorMsgs().KeyValuePair, "Owned", c.Owned)
		}
		if column == "ID" {
			fmt.Printf(ErrorMsgs().KeyValuePair, "updating Joined of", c.Name)
			c.Joined = true
			fmt.Printf(ErrorMsgs().KeyValuePair, "Joined", c.Joined)
		}
		channels = append(channels, c)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}
	if column == "OwnerID" {
		fmt.Printf(ErrorMsgs().KeyValuePair, "Channels owned by current user", len(channels))
	}
	return channels, nil
}

func (m *ChannelModel) All() ([]models.Channel, error) {
	stmt := "SELECT ID, OwnerID, Name, Avatar, Banner, Description, Created, Privacy, IsFlagged, IsMuted FROM Channels ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
		}
	}()

	var Channels []models.Channel
	for rows.Next() {
		p := models.Channel{}
		err = rows.Scan(&p.ID, &p.OwnerID, &p.Name, &p.Avatar, &p.Banner, &p.Description, &p.Created, &p.Privacy, &p.IsFlagged, &p.IsMuted)
		if err != nil {
			return nil, err
		}
		Channels = append(Channels, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	fmt.Printf(ErrorMsgs().KeyValuePair, "Total channels", len(Channels))
	return Channels, nil
}

// Search queries the database for any channel column that contains the value and returns a slice of matching channels
func (m *ChannelModel) Search(column string, value interface{}) ([]models.Channel, error) {
	// Validate column name to prevent SQL injection
	validColumns := map[string]bool{
		"id":          true,
		"ownerId":     true,
		"name":        true,
		"avatar":      true,
		"banner":      true,
		"description": true,
		"created":     true,
		"privacy":     true,
		"isMuted":     true,
		"isFlagged":   true,
	}

	if !validColumns[column] {
		return nil, fmt.Errorf("invalid column name: %s", column)
	}

	// Base query
	query := fmt.Sprintf("SELECT id, ownerId, name, avatar, banner, description, created, privacy, isMuted, isFlagged FROM channels WHERE %s = ?", column)

	// Execute the query
	rows, err := m.DB.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse results
	var channels []models.Channel
	for rows.Next() {
		var channel models.Channel
		var avatar, banner sql.NullString

		if err := rows.Scan(
			&channel.ID,
			&channel.OwnerID,
			&channel.Name,
			&avatar,
			&banner,
			&channel.Description,
			&channel.Created,
			&channel.Privacy,
			&channel.IsMuted,
			&channel.IsFlagged,
		); err != nil {
			return nil, err
		}

		channel.Avatar = avatar.String
		channel.Banner = banner.String
		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return channels, nil
}

func (m *ChannelModel) AddPostToChannel(channelID, postID int) error {
	stmt := "INSERT INTO PostChannels (ChannelID, PostID, Created) VALUES (?, ?, DateTime('now'))"
	_, err := m.DB.Exec(stmt, channelID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (m *ChannelModel) GetPostsFromChannel(channelID int) error {
	stmt := "SELECT PostID FROM PostChannels WHERE ChannelID = ?"
	_, err := m.DB.Exec(stmt, channelID)
	if err != nil {
		return err
	}
	return nil
}
