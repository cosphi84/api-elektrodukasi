package repositories

import (
	"database/sql"
	"errors"

	"elektrod/internal/models"

	"github.com/google/uuid"
	nestedset "github.com/longbridgeapp/nested-set"
	"gorm.io/gorm"
)

type CommentRepository interface {
	FindAll(articleID uuid.UUID) ([]models.Comment, error)
	FindByID(id int64) (*models.Comment, error)
	Create(comment *models.Comment, parent *models.Comment) error
	Update(comment *models.Comment) error
	SoftDelete(id int64) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

// FindAll returns the full comment tree for an article, ordered by lft ASC.
// lft ASC gives natural top-down tree order — parent always before children.
func (r *commentRepository) FindAll(articleID uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.
		Preload("User").
		Where("article_id = ? AND deleted_at IS NULL", articleID).
		Order("lft ASC").
		Find(&comments).Error
	return comments, err
}

// FindByID fetches a single comment — used before Update/Delete/Reply.
func (r *commentRepository) FindByID(id int64) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.
		Preload("User").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	return &comment, nil
}

// Create inserts a new comment node into the nested set tree.
// parent == nil → root level comment on the article.
// parent != nil → reply to that comment.
func (r *commentRepository) Create(comment *models.Comment, parent *models.Comment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return nestedset.Create(tx, comment, parent)
	})
}

// Update modifies content and approved flag only.
// Never touch nested set fields (lft, rgt, depth) here — use MoveTo for that.
func (r *commentRepository) Update(comment *models.Comment) error {
	return r.db.Model(comment).
		Select("content", "approved", "updated_at").
		Updates(comment).Error
}

// SoftDelete marks a comment as deleted without removing it from the tree.
// The node stays in the nested set — orphaning children is intentional
// (show "deleted comment" placeholder in UI instead of collapsing the tree).
func (r *commentRepository) SoftDelete(id int64) error {
	result := r.db.Model(&models.Comment{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCommentNotFound
	}
	return nil
}

// --- Helpers ---

// BuildParentID converts *int64 to sql.NullInt64 for the nested set model.
// Call this before Create when mapping from DTO.
func BuildParentID(parentID *int64) sql.NullInt64 {
	if parentID == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *parentID, Valid: true}
}
