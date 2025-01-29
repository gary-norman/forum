package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type MembershipModel struct {
	DB *sql.DB
}

func (m *MembershipModel) Insert(userID, channelID int) error {
	stmt := "INSERT INTO Memberships (UserID, ChannelID, Created) VALUES (?, ?, DateTime('now'))"
	_, err := m.DB.Exec(stmt, userID, channelID)
	return err
}

func (m *MembershipModel) UserMemberships(userID int) ([]models.Membership, error) {
	stmt := "SELECT ID, UserID, ChannelID, Created FROM Memberships WHERE ID = ?"
	rows, queryErr := m.DB.Query(stmt, userID)
	if queryErr != nil {
		return nil, queryErr
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "UserMemberships", closeErr)
		}
	}()
	var memberships []models.Membership
	for rows.Next() {
		p := models.Membership{}
		scanErr := rows.Scan(&p.ID, &p.UserID, &p.ChannelID, &p.Created)
		if scanErr != nil {
			return nil, scanErr
		}
		memberships = append(memberships, p)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}
	fmt.Printf(ErrorMsgs().KeyValuePair, "Channels", len(memberships))
	return memberships, nil
}

func (m *MembershipModel) All() ([]models.Membership, error) {
	stmt := "SELECT ID, UserID, ChannelID, Created FROM Memberships ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
		}
	}()

	var Memberships []models.Membership
	for rows.Next() {
		p := models.Membership{}
		err = rows.Scan(&p.ID, &p.UserID, &p.ChannelID, &p.Created)
		if err != nil {
			return nil, err
		}
		Memberships = append(Memberships, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Memberships, nil
}
