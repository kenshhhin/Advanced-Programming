package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type News struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}

type Comment struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	NewsID int    `json:"news_id"`
	Text   string `json:"text"`
}
