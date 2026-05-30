package dto

import "github.com/google/uuid"

// --- Category ---

type CreateCategoryRequest struct {
	Name        string  `json:"name"        validate:"required,min=2,max=100"`
	Description *string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name"        validate:"omitempty,min=2,max=100"`
	Description *string `json:"description"`
}

// --- Tag ---

type CreateTagRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type UpdateTagRequest struct {
	Name *string `json:"name" validate:"omitempty,min=2,max=50"`
}

// --- Project ---

type CreateProjectRequest struct {
	Title    string     `json:"title"    validate:"required,min=2,max=255"`
	Slug     string     `json:"slug"     validate:"required"`
	Summary  *string    `json:"summary"`
	Link     *string    `json:"link"     validate:"omitempty,url"`
	Metadata []byte     `json:"metadata"` // raw JSONB
	OwnerID  *uuid.UUID `json:"owner_id"`
}

type UpdateProjectRequest struct {
	Title    *string    `json:"title"    validate:"omitempty,min=2,max=255"`
	Slug     *string    `json:"slug"`
	Summary  *string    `json:"summary"`
	Link     *string    `json:"link"     validate:"omitempty,url"`
	Metadata []byte     `json:"metadata"`
	OwnerID  *uuid.UUID `json:"owner_id"`
}

// --- User ---

type CreateUserRequest struct {
	Name     string  `json:"name"     validate:"required,min=2,max=100"`
	Email    string  `json:"email"    validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
	Role     string  `json:"role"     validate:"omitempty,oneof=user admin"`
	Avatar   *string `json:"avatar"`
}

type UpdateUserRequest struct {
	Name   *string `json:"name"   validate:"omitempty,min=2,max=100"`
	Email  *string `json:"email"  validate:"omitempty,email"`
	Avatar *string `json:"avatar"`
	Role   *string `json:"role"   validate:"omitempty,oneof=user admin"`
}

// --- Comment ---

type CreateCommentRequest struct {
	ArticleID uuid.UUID  `json:"article_id" validate:"required"`
	UserID    *uuid.UUID `json:"user_id"`
	ParentID  *int64     `json:"parent_id"` // nil = root comment
	Content   string     `json:"content"    validate:"required,min=1"`
}

type UpdateCommentRequest struct {
	Content  *string `json:"content"  validate:"omitempty,min=1"`
	Approved *bool   `json:"approved"`
}
