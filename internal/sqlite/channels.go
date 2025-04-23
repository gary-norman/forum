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

func (m *ChannelModel) Insert(ownerID int64, name, description, avatar, banner string, privacy, isFlagged, isMuted bool) error {
	stmt := "INSERT INTO Channels (OwnerID, Name, Description, Created, Avatar, Banner, Privacy, IsFlagged, IsMuted) VALUES (?, ?, ?, DateTime('now'), ?, ?, ?, ?, ?)"
	_, err := m.DB.Exec(stmt, ownerID, name, description, avatar, banner, privacy, isFlagged, isMuted)
	return err
}

func (m *ChannelModel) OwnedOrJoinedByCurrentUser(ID int64, column string) ([]models.Channel, error) {
	// Validate column name to prevent SQL injection
	if !isValidColumn(column) {
		return nil, fmt.Errorf("invalid column name provided: %s", column)
	}

	// Base query
	query := "SELECT id, ownerId, name, avatar, banner, description, created, privacy, isMuted, isFlagged FROM channels WHERE " + column + " = ?"

	// Execute the query
	rows, err := m.DB.Query(query, ID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Parse results
	channels := make([]models.Channel, 0) // Pre-allocate slice
	for rows.Next() {
		c, err := parseChannelRow(rows)
		if err != nil {
			return nil, fmt.Errorf("error parsing row: %w", err)
		}
		if column == "OwnerID" {
			// fmt.Printf(ErrorMsgs().KeyValuePair, "updating Owned of", c.Name)
			c.Owned = true
			// fmt.Printf(ErrorMsgs().KeyValuePair, "Owned", c.Owned)
		}
		if column == "ID" {
			// fmt.Printf(ErrorMsgs().KeyValuePair, "updating Joined of", c.Name)
			c.Joined = true
			// fmt.Printf(ErrorMsgs().KeyValuePair, "Joined", c.Joined)
		}
		channels = append(channels, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	if column == "OwnerID" {
		// fmt.Printf(ErrorMsgs().KeyValuePair, "Channels owned by current user", len(channels))
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
	// fmt.Printf(ErrorMsgs().KeyValuePair, "Total channels", len(Channels))
	return Channels, nil
}

// SearchChannelsByColumn queries the database for any channel column that contains the value and returns a slice of matching channels
func (m *ChannelModel) SearchChannelsByColumn(column string, value interface{}) ([]models.Channel, error) {
	// Validate column name to prevent SQL injection
	if !isValidColumn(column) {
		return nil, fmt.Errorf("invalid column name provided: %s", column)
	}

	// Base query
	query := "SELECT ID, OwnerId, Name, Avatar, Banner, Description, Created, Privacy, IsMuted, IsFlagged FROM channels WHERE " + column + " = ?"

	// Execute the query
	rows, err := m.DB.Query(query, value)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Parse results
	channels := make([]models.Channel, 0) // Pre-allocate slice
	for rows.Next() {
		channel, err := parseChannelRow(rows)
		if err != nil {
			return nil, fmt.Errorf("error parsing row: %w", err)
		}
		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return channels, nil
}

func isValidColumn(column string) bool {
	validColumns := map[string]bool{
		"ID":          true,
		"OwnerID":     true,
		"Name":        true,
		"Avatar":      true,
		"Banner":      true,
		"Description": true,
		"Created":     true,
		"Privacy":     true,
		"IsMuted":     true,
		"IsFlagged":   true,
	}
	return validColumns[column]
}

func parseChannelRow(rows *sql.Rows) (models.Channel, error) {
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
		return channel, err
	}

	channel.Avatar = avatar.String
	channel.Banner = banner.String
	return channel, nil
}

func (m *ChannelModel) AddPostToChannel(channelID, postID int64) error {
	stmt := "INSERT INTO PostChannels (ChannelID, PostID, Created) VALUES (?, ?, DateTime('now'))"
	_, err := m.DB.Exec(stmt, channelID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (m *ChannelModel) GetPostIDsFromChannel(channelID int64) ([]int64, error) {
	var postIDs []int64
	stmt := "SELECT PostID FROM PostChannels WHERE ChannelID = ?"
	rows, err := m.DB.Query(stmt, channelID)
	if err != nil {
		return postIDs, err
	}
	defer rows.Close()

	for rows.Next() {
		var postID int64
		if err := rows.Scan(&postID); err != nil {
			return postIDs, err
		}
		postIDs = append(postIDs, postID)
	}

	return postIDs, nil
}

func (m *ChannelModel) GetChannelIdFromPost(postID int64) ([]int64, error) {
	var channelIDs []int64
	stmt := "SELECT ChannelID FROM PostChannels WHERE PostID = ?"
	rows, err := m.DB.Query(stmt, postID)
	if err != nil {
		return channelIDs, err
	}
	defer rows.Close()

	for rows.Next() {
		var channelID int64
		if err := rows.Scan(&channelID); err != nil {
			return channelIDs, err
		}
		channelIDs = append(channelIDs, channelID)
	}

	return channelIDs, nil
}

func (m *ChannelModel) GetChannelNameFromID(id int64) (string, error) {
	var name string
	stmt := "SELECT Name FROM Channels WHERE ID = ?"
	row := m.DB.QueryRow(stmt, id)
	if err := row.Scan(&name); err != nil {
		return "", err
	}

	return name, nil
}
