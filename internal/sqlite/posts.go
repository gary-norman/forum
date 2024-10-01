package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content string) error {
	stmt := "INSERT INTO posts (title, content, created_at) VALUES (?, ?, DateTime('now'))"
	_, err := m.DB.Exec(stmt, title, content)
	return err
}

func (m *PostModel) All() ([]models.Posts, error) {
	ErrorMsgs := models.CreateErrorMessages()
	//TODO these will need to be changes according to the new DB
	stmt := "SELECT post_id, thread_id, user_id, content, created_at, title FROM posts ORDER BY post_id DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	var posts []models.Posts
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf(ErrorMsgs.Close, "rows", "All")
		}
	}(rows)
	for rows.Next() {
		p := models.Posts{}
		err := rows.Scan(&p.ID, &p.ThreadId, &p.UserId, &p.Content, &p.Created, &p.Title)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}
