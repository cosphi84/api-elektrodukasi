package models

import "time"

type Tag struct {
	ID        int       `json:"id"`
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
}
