package repositories

import (
	"elektrod/internal/dto"
	"elektrod/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- Interface ---

type ArticleRepository interface {
	List(query dto.ArticleListQuery) (dto.PaginatedResult[models.Article], error)
	Create(article *models.Article, tagIDs []uuid.UUID) error
	FindByID(id uuid.UUID) (*models.Article, error)
	Update(article *models.Article, tagIDs []uuid.UUID) error
	SoftDelete(id uuid.UUID) error
}

// --- Implementation ---

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// List returns paginated articles sorted by newest first.
// Supports optional published filter and title search.
func (r *articleRepository) List(query dto.ArticleListQuery) (dto.PaginatedResult[models.Article], error) {
	query.Normalize()

	tx := r.db.Model(&models.Article{}).
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Where("deleted_at IS NULL").
		Order("created_at DESC")

	// Optional: filter by published status
	if query.Published != nil {
		tx = tx.Where("published = ?", *query.Published)
	}

	// Optional: search by title (case-insensitive)
	if query.Search != "" {
		tx = tx.Where("title ILIKE ?", "%"+query.Search+"%")
	}

	// Count total before applying limit/offset
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return dto.PaginatedResult[models.Article]{}, err
	}

	// Fetch paginated data
	var articles []models.Article
	if err := tx.
		Limit(query.PerPage).
		Offset(query.Offset()).
		Find(&articles).Error; err != nil {
		return dto.PaginatedResult[models.Article]{}, err
	}

	return dto.NewPaginatedResult(articles, total, query.PaginationQuery), nil
}

// Create inserts a new article and associates tags in a single transaction.
func (r *articleRepository) Create(article *models.Article, tagIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		if len(tagIDs) > 0 {
			tags, err := findTagsByIDs(tx, tagIDs)
			if err != nil {
				return err
			}
			if err := tx.Model(article).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}

		return nil
	})
}

// FindByID fetches a single article by ID for editing purposes.
// Returns ErrArticleNotFound if not found or already soft-deleted.
func (r *articleRepository) FindByID(id uuid.UUID) (*models.Article, error) {
	var article models.Article
	err := r.db.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&article).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}

	return &article, nil
}

// Update modifies an existing article and replaces its tag associations.
func (r *articleRepository) Update(article *models.Article, tagIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Save all fields except created_at and deleted_at
		if err := tx.Model(article).
			Omit("created_at", "deleted_at").
			Save(article).Error; err != nil {
			return err
		}

		// Replace tags (clear + re-associate)
		tags, err := findTagsByIDs(tx, tagIDs)
		if err != nil {
			return err
		}
		if err := tx.Model(article).Association("Tags").Replace(tags); err != nil {
			return err
		}

		return nil
	})
}

// SoftDelete sets deleted_at timestamp without removing the row.
func (r *articleRepository) SoftDelete(id uuid.UUID) error {
	result := r.db.Model(&models.Article{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrArticleNotFound // already deleted or never existed
	}

	return nil
}

// --- Helpers ---

// findTagsByIDs loads Tag models from a list of UUIDs.
// Fails fast if any ID doesn't exist — no silent partial associations.
func findTagsByIDs(tx *gorm.DB, ids []uuid.UUID) ([]models.Tag, error) {
	if len(ids) == 0 {
		return []models.Tag{}, nil
	}

	var tags []models.Tag
	if err := tx.Where("id IN ?", ids).Find(&tags).Error; err != nil {
		return nil, err
	}

	if len(tags) != len(ids) {
		return nil, ErrTagNotFound
	}

	return tags, nil
}
