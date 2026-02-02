package models

import "time"

type Author struct {
	ID        string    `json:"uuid"`
	Name      string    `json:"name"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Blocked   bool      `json:"blocked"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
