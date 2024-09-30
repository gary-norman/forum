package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) All() ([]models.Post, error) {
	stmt := "SELECT user_id, title, content, created_at FROM posts ORDER BY post_id DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	posts := []models.Post{}
	for rows.Next() {
		p := models.Post{}
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
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
