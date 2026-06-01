package services

import (
	"errors"
	"elektrod/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentService interface {
	CreateComment(comment *models.Comment) error
	GetCommentByID(id int64) (*models.Comment, error)
	GetCommentsByArticle(articleID uuid.UUID, page, pageSize int) ([]models.Comment, int64, error)
	GetCommentsByArticleApproved(articleID uuid.UUID, page, pageSize int) ([]models.Comment, int64, error)
	UpdateComment(id int64, updates map[string]interface{}) (*models.Comment, error)
	DeleteComment(id int64) error
	ApproveComment(id int64) error
	GetCommentReplies(parentID int64) ([]models.Comment, error)
}

type commentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) CommentService {
	return &commentService{db: db}
}

func (s *commentService) CreateComment(comment *models.Comment) error {
	return s.db.Create(comment).Error
}

func (s *commentService) GetCommentByID(id int64) (*models.Comment, error) {
	var comment models.Comment
	err := s.db.Preload("User").Preload("Article").
		Where("id = ?", id).First(&comment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("comment not found")
	}
	return &comment, err
}

func (s *commentService) GetCommentsByArticle(articleID uuid.UUID, page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	offset := (page - 1) * pageSize
	err := s.db.Model(&models.Comment{}).Where("article_id = ?", articleID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Where("article_id = ?", articleID).
		Preload("User").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error
	return comments, total, err
}

func (s *commentService) GetCommentsByArticleApproved(articleID uuid.UUID, page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	offset := (page - 1) * pageSize
	err := s.db.Model(&models.Comment{}).
		Where("article_id = ? AND approved = ?", articleID, true).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Where("article_id = ? AND approved = ?", articleID, true).
		Preload("User").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error
	return comments, total, err
}

func (s *commentService) UpdateComment(id int64, updates map[string]interface{}) (*models.Comment, error) {
	comment, err := s.GetCommentByID(id)
	if err != nil {
		return nil, err
	}

	err = s.db.Model(comment).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) DeleteComment(id int64) error {
	return s.db.Where("id = ?", id).Delete(&models.Comment{}).Error
}

func (s *commentService) ApproveComment(id int64) error {
	return s.db.Model(&models.Comment{}).Where("id = ?", id).Update("approved", true).Error
}

func (s *commentService) GetCommentReplies(parentID int64) ([]models.Comment, error) {
	var comments []models.Comment
	err := s.db.Where("parent_id = ?", parentID).
		Preload("User").
		Order("created_at ASC").Find(&comments).Error
	return comments, err
}
