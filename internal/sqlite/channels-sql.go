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

func (m *ChannelModel) Insert(ownerID models.UUIDField, name, description, avatar, banner string, privacy, isFlagged, isMuted bool) error {
	stmt := "INSERT INTO Channels (OwnerID, Name, Description, Created, Avatar, Banner, Privacy, IsFlagged, IsMuted) VALUES (?, ?, ?, DateTime('now'), ?, ?, ?, ?, ?)"
	_, err := m.DB.Exec(stmt, ownerID, name, description, avatar, banner, privacy, isFlagged, isMuted)
	return err
}

func (m *ChannelModel) OwnedOrJoinedByCurrentUser(ID models.UUIDField) ([]models.Channel, error) {
	stmt := `
	SELECT c.*,
	COUNT(m.UserID) AS MemberCount
	From Channels c
	LEFT JOIN Memberships m ON c.ID = m.ChannelID
	WHERE c.ID IN (
		SELECT ChannelID FROM Memberships WHERE UserID = ?
	)
	OR c.OwnerID = ?
	GROUP BY c.ID
	ORDER BY Name DESC
	`
	rows, err := m.DB.Query(stmt, ID, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse results
	channels := make([]models.Channel, 0) // Pre-allocate slice
	for rows.Next() {
		c, err := parseChannelRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error parsing row: %w", err)
		}
		// FIXME: This is a temporary fix to set the channel as joined:we need to come up with a more robust solution
		c.Joined = true
		// TODO (realtime) get this data freom websockets
		c.MembersOnline = 0
		channels = append(channels, *c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return channels, nil
}

func (m *ChannelModel) IsUserMemberOfChannel(userID models.UUIDField, channelID int64) (bool, error) {
	var exists int
	stmt := `
		SELECT EXISTS (
			SELECT 1 FROM Memberships
			WHERE UserID = ? AND ChannelID = ?
		)
	`
	err := m.DB.QueryRow(stmt, userID, channelID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (m *ChannelModel) GetChannelsByID(id int64) ([]models.Channel, error) {
	stmt := `
	SELECT c.*, 
  COUNT(m.UserID) AS MemberCount
	FROM Channels c
	LEFT JOIN Memberships m ON c.ID = m.ChannelID
	WHERE c.ID = ?
	GROUP BY c.ID;
	`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse results
	channels := make([]models.Channel, 0) // Pre-allocate slice
	for rows.Next() {
		c, err := parseChannelRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error parsing row: %w", err)
		}
		// TODO (realtime) get this data freom websockets
		c.MembersOnline = 0
		channels = append(channels, *c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return channels, nil
}

func (m *ChannelModel) GetChannelByID(id int64) (*models.Channel, error) {
	stmt := `
	SELECT c.*, 
  COUNT(m.UserID) AS MemberCount
	FROM Channels c
	LEFT JOIN Memberships m ON c.ID = m.ChannelID
	WHERE c.ID = ?
	GROUP BY c.ID;
	`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	channels := make([]models.Channel, 0) // Pre-allocate slice
	for rows.Next() {
		c, err := parseChannelRows(rows)
		if err != nil {
			return nil, err
		}
		// TODO (realtime) get this data freom websockets
		c.MembersOnline = 0
		channels = append(channels, *c)
	}
	return &channels[0], nil
}

func (m *ChannelModel) GetNameOfChannel(channelID int64) (string, error) {
	stmt := "SELECT Name FROM Channels WHERE ID = ?)"
	rows, err := m.DB.Query(stmt, channelID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var username string
	for rows.Next() {
		if err := rows.Scan(&username); err != nil {
			return "", err
		}
	}
	return username, nil
}

func (m *ChannelModel) GetNameOfChannelOwner(channelID int64) (string, error) {
	stmt := "SELECT Username FROM Users WHERE ID = (SELECT OwnerID FROM Channels WHERE ID = ?)"
	rows, err := m.DB.Query(stmt, channelID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var username string
	for rows.Next() {
		if err := rows.Scan(&username); err != nil {
			return "", err
		}
	}
	return username, nil
}

func (m *ChannelModel) All() ([]models.Channel, error) {
	stmt := `
	SELECT c.*, 
  COUNT(m.UserID) AS MemberCount
	FROM Channels c
	LEFT JOIN Memberships m ON c.ID = m.ChannelID
	GROUP BY c.ID;
	`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
		}
	}()

	channels := make([]models.Channel, 0) // Pre-allocate slice
	for rows.Next() {
		c, err := parseChannelRows(rows)
		if err != nil {
			return nil, err
		}
		// TODO (realtime) get this data freom websockets
		c.MembersOnline = 0
		channels = append(channels, *c)
	}
	// fmt.Printf(ErrorMsgs().KeyValuePair, "Total channels", len(Channels))
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

func parseChannelRow(row *sql.Row) (*models.Channel, error) {
	var channel models.Channel
	var avatar, banner sql.NullString

	if err := row.Scan(
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
	models.UpdateTimeSince(&channel)
	return &channel, nil
}

func parseChannelRows(rows *sql.Rows) (*models.Channel, error) {
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
		&channel.Members,
	); err != nil {
		return nil, err
	}

	channel.Avatar = avatar.String
	channel.Banner = banner.String
	models.UpdateTimeSince(&channel)
	return &channel, nil
}
