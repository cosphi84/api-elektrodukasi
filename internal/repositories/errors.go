package repositories

import "errors"

var (
	ErrArticleNotFound  = errors.New("article not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrProjectNotFound  = errors.New("project not found")
	ErrCommentNotFound  = errors.New("comment not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrTagNotFound      = errors.New("one or more tags not found")
)
