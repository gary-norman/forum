package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gary-norman/forum/internal/models"
)

type RuleModel struct {
	DB *sql.DB
}

// CreateRule inserts a new rule into the Rules table
func (m *RuleModel) CreateRule(rule string) (int64, error) {
	var ruleID int64
	stmtRule := "INSERT INTO Rules (Rule, Created, Predefined) VALUES (?, DateTime('now'), 0)"
	result, err := m.DB.Exec(stmtRule, rule)
	if err != nil {
		return 0, err
	}

	// Get the last inserted ID
	ruleID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return ruleID, nil
}

// InsertRule inserts a rule:channel reference into the ChannelsRules table
func (m *RuleModel) InsertRule(channelID, ruleID int64) error {
	stmt := "INSERT INTO ChannelsRules (ChannelID, RuleID) VALUES (?, ?)"
	_, channelRuleErr := m.DB.Exec(stmt, channelID, ruleID)
	return channelRuleErr
}

// InsertChannelRule adds an existing rule to the ChannelsRules table, omitting if it already exists
func (m *RuleModel) InsertChannelRule(channelID, ruleID int64) error {
	stmt := "INSERT INTO ChannelsRules (ChannelID, RuleID) VALUES (?, ?) ON CONFLICT(ChannelID, RuleID) DO NOTHING"
	_, channelRuleErr := m.DB.Exec(stmt, channelID, ruleID)
	return channelRuleErr
}

// EditRule edits the rule string in the Rules table
func (m *RuleModel) EditRule(id int64, rule string) error {
	stmt := "UPDATE Rules SET Rule = ? WHERE ID = ?"
	_, editErr := m.DB.Exec(stmt, rule, id)
	return editErr
}

// DeleteRule removes a rule/channel reference from the ChannelsRules table
func (m *RuleModel) DeleteRule(channelID, ruleID int64) error {
	stmt := "DELETE FROM ChannelsRules WHERE ChannelID = ? AND RuleID = ?"
	_, deleteErr := m.DB.Exec(stmt, channelID, ruleID)
	return deleteErr
}

// All returns every row from the Rules table ordered by ID, descending
func (m *RuleModel) All() ([]models.Rule, error) {
	stmt := "SELECT * FROM Rules ORDER BY ID DESC"
	rows, queryErr := m.DB.Query(stmt)
	if queryErr != nil {
		return nil, fmt.Errorf(ErrorMsgs.Query, "Rules", queryErr)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, "rows", "All")
		}
	}()

	var Rules []models.Rule
	for rows.Next() {
		r := models.Rule{}
		scanErr := rows.Scan(&r.ID, &r.Rule, &r.Created, &r.Predefined)
		if scanErr != nil {
			return nil, scanErr
		}
		Rules = append(Rules, r)
	}

	if queryErr = rows.Err(); queryErr != nil {
		return nil, fmt.Errorf(ErrorMsgs.Query, "Rules", queryErr)
	}

	return Rules, nil
}

func (m *RuleModel) AllForChannel(channelID int64) ([]models.Rule, error) {
	// fetch the references from ChannelsRules
	refStmt := "SELECT RuleID FROM ChannelsRules WHERE ChannelID = ?"
	crRows, queryErr := m.DB.Query(refStmt, channelID)
	if queryErr != nil {
		return nil, fmt.Errorf(ErrorMsgs.Query, "AllForChannel", queryErr)
	}
	defer func() {
		if closeErr := crRows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, "rows", "All")
		}
	}()
	var IDs []int
	for crRows.Next() {
		var i int
		scanErr := crRows.Scan(&i)
		if scanErr != nil {
			return nil, scanErr
		}
		IDs = append(IDs, i)
	}

	// prepare the statement for use in the loop
	ruleStmt, insertErr := m.DB.Prepare("SELECT * FROM Rules WHERE ID = ?")
	if insertErr != nil {
		return nil, insertErr
	}
	defer func(stmt *sql.Stmt) {
		closErr := stmt.Close()
		if closErr != nil {
			log.Printf(ErrorMsgs.Close, "stmt", "insert", closErr)
		}
	}(ruleStmt)

	var Rules []models.Rule
	// create a []rule from the slice of rule IDs
	for _, ruleID := range IDs {
		rows, err := ruleStmt.Query(ruleID)
		if err != nil {
			return Rules, err
		}
		for rows.Next() {
			r := models.Rule{}
			scanErr := rows.Scan(&r.ID, &r.Rule, &r.Created, &r.Predefined)
			if scanErr != nil {
				return nil, scanErr
			}
			Rules = append(Rules, r)
		}
	}
	return Rules, nil
}
