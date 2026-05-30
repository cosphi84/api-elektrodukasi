package dto

import (
	"github.com/google/uuid"
)

// --- Query ---

type ArticleListQuery struct {
	PaginationQuery
	Published *bool  `form:"published" json:"published"` // nil = all, true/false = filter
	Search    string `form:"search"    json:"search"`    // search by title
}

// --- Request ---

type CreateArticleRequest struct {
	Title       string      `json:"title"        validate:"required,min=3,max=255"`
	Slug        string      `json:"slug"         validate:"required,slug"`
	Summary     *string     `json:"summary"`
	ContentHTML *string     `json:"content_html"`
	ContentJSON []byte      `json:"content_json"` // raw tiptap JSON
	Image       *string     `json:"image"`
	CategoryID  *uuid.UUID  `json:"category_id"`
	TagIDs      []uuid.UUID `json:"tag_ids"`
	Published   bool        `json:"published"`
}

type UpdateArticleRequest struct {
	Title       *string     `json:"title"        validate:"omitempty,min=3,max=255"`
	Slug        *string     `json:"slug"         validate:"omitempty,slug"`
	Summary     *string     `json:"summary"`
	ContentHTML *string     `json:"content_html"`
	ContentJSON []byte      `json:"content_json"`
	Image       *string     `json:"image"`
	CategoryID  *uuid.UUID  `json:"category_id"`
	TagIDs      []uuid.UUID `json:"tag_ids"`
	Published   *bool       `json:"published"`
}

// --- Response ---

type ArticleResponse struct {
	ID          uuid.UUID         `json:"id"`
	Title       string            `json:"title"`
	Slug        string            `json:"slug"`
	Summary     *string           `json:"summary,omitempty"`
	ContentHTML *string           `json:"content_html,omitempty"`
	ContentJSON []byte            `json:"content_json,omitempty"`
	Image       *string           `json:"image,omitempty"`
	Published   bool              `json:"published"`
	PublishedAt *string           `json:"published_at,omitempty"`
	ViewCount   int64             `json:"view_count"`
	CreatedAt   string            `json:"created_at"`
	Author      *AuthorResponse   `json:"author,omitempty"`
	Category    *CategoryResponse `json:"category,omitempty"`
	Tags        []TagResponse     `json:"tags,omitempty"`
}

type AuthorResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Avatar *string   `json:"avatar,omitempty"`
}

type CategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type TagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
