package models

import "time"

type Article struct {
	ID         string    `json:"uuid"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Content    string    `json:"content"`
	Tags       string    `json:"tags"`
	Cover      string    `json:"cover"`
	Published  bool      `json:"published"`
	Featured   bool      `json:"featured"`
	Hit        int       `json:"hit"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy string    `json:"modified_by"`
	Category   string    `json:"category"`
	CategoryId string    `json:"category_id"`
	Author     string    `json:"author"`
	AuthorId   string    `json:"author_id"`
}
