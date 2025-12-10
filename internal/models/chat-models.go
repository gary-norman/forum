package models

import "time"

type Chat struct {
	ID         UUIDField     `json:"id"`
	ChatType   string        `json:"type"`
	Name       string        `json:"name"`
	Created    time.Time     `json:"created"`
	LastActive time.Time     `json:"last_active"`
	Group      Group         `json:"group"`
	Buddy      User          `json:"buddy"`
	Messages   []ChatMessage `json:"messages"`
}

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ChatMessage struct {
	ID      UUIDField `json:"id"`
	ChatID  UUIDField `json:"chat_id"`
	UserID  UUIDField `json:"user_id"`
	Created time.Time `json:"created"`
	Content string    `json:"content"`
}
