package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type RuleModel struct {
	DB *sql.DB
}

// CreateRule inserts a new rule into the Rules table
func (m *RuleModel) CreateRule(rule string) (int, error) {
	var ruleId int
	stmtRule := "INSERT INTO Rules (Rule, Created, Predefined) VALUES (?, DateTime('now'), 0)"
	result, ruleErr := m.DB.Exec(stmtRule, rule)
	if ruleErr != nil {
		return 0, ruleErr
	}

	// Get the last inserted ID
	ruleId64, idErr := result.LastInsertId()
	if idErr != nil {
		return 0, idErr
	}

	ruleId = int(ruleId64)
	return ruleId, nil
}

// InsertRule inserts a rule:channel reference into the ChannelsRules table
func (m *RuleModel) InsertRule(channelId, ruleId int) error {
	stmt := "INSERT INTO ChannelsRules (ChannelID, RuleID) VALUES (?, ?)"
	_, channelRuleErr := m.DB.Exec(stmt, channelId, ruleId)
	return channelRuleErr
}

// InsertChannelRule adds an existing rule to the ChannelsRules table, omitting if it already exists
func (m *RuleModel) InsertChannelRule(channelId, ruleId int) error {
	stmt := "INSERT INTO ChannelsRules (ChannelID, RuleID) VALUES (?, ?) ON CONFLICT(ChannelID, RuleID) DO NOTHING"
	_, channelRuleErr := m.DB.Exec(stmt, channelId, ruleId)
	return channelRuleErr
}

// EditRule edits the rule string in the Rules table
func (m *RuleModel) EditRule(id int, rule string) error {
	stmt := "UPDATE Rules SET Rule = ? WHERE ID = ?"
	_, editErr := m.DB.Exec(stmt, rule, id)
	return editErr
}

// DeleteRule removes a rule/channel reference from the ChannelsRules table
func (m *RuleModel) DeleteRule(channelId, ruleId int) error {
	stmt := "DELETE FROM ChannelsRules WHERE ChannelID = ? AND RuleID = ?"
	_, deleteErr := m.DB.Exec(stmt, channelId, ruleId)
	return deleteErr
}

// All returns every row from the Rules table ordered by ID, descending
func (m *RuleModel) All() ([]models.Rule, error) {
	stmt := "SELECT * FROM Rules ORDER BY ID DESC"
	rows, queryErr := m.DB.Query(stmt)
	if queryErr != nil {
		return nil, errors.New(fmt.Sprintf(ErrorMsgs().Query, "Rules", queryErr))
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, "rows", "All")
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
		return nil, errors.New(fmt.Sprintf(ErrorMsgs().Query, "Rules", queryErr))
	}

	return Rules, nil
}

func (m *RuleModel) AllForChannel(channelId int) ([]models.Rule, error) {
	//fetch the references from ChannelsRules
	refStmt := "SELECT RuleID FROM ChannelsRules WHERE ChannelID = ?"
	crRows, queryErr := m.DB.Query(refStmt, channelId)
	if queryErr != nil {
		return nil, errors.New(fmt.Sprintf(ErrorMsgs().Query, "AllForChannel", queryErr))
	}
	defer func() {
		if closeErr := crRows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, "rows", "All")
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
			log.Printf(ErrorMsgs().Close, "stmt", "insert", closErr)
		}
	}(ruleStmt)

	var Rules []models.Rule
	// create a []rule from the slice of rule IDs
	for _, ruleId := range IDs {
		rows, err := ruleStmt.Query(ruleId)
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
